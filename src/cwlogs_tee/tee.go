package cwlogs_tee

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"
	"github.com/cenkalti/backoff"
	"io"
	"regexp"
	"strings"
	"time"
)

func backoffRateExceeded(operation func() error) (err error) {
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	for _ = range ticker.C {
		err = operation()

		if err == nil {
			ticker.Stop()
			break
		}

		awsErr, ok := err.(awserr.Error)

		if !ok {
			ticker.Stop()
			break
		}

		if awsErr.Code() != "ThrottlingException" {
			ticker.Stop()
			break
		}
	}

	return err
}

type CWLogsTee struct {
	LogGroupName  string
	LogStreamName string
	In            io.Reader
	Out           io.Writer
	Now           func() time.Time
}

func (tee *CWLogsTee) scan(fn func(string) error) (err error) {
	scanner := bufio.NewScanner(tee.In)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimRight(line, "\n")
		err = fn(line)

		if err != nil {
			return
		}
	}

	err = scanner.Err()

	return
}

func (tee *CWLogsTee) isGroupExist(svc cloudwatchlogsiface.CloudWatchLogsAPI) (exist bool, err error) {
	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(tee.LogGroupName),
	}

	err = backoffRateExceeded(func() error {
		err = svc.DescribeLogGroupsPages(params, func(page *cloudwatchlogs.DescribeLogGroupsOutput, lastPage bool) bool {
			for _, group := range page.LogGroups {
				if *group.LogGroupName == tee.LogGroupName {
					exist = true
					return false
				}
			}

			return !lastPage
		})

		return err
	})

	return
}

func (tee *CWLogsTee) createLogGroup(svc cloudwatchlogsiface.CloudWatchLogsAPI) (err error) {
	params := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(tee.LogGroupName),
	}

	err = backoffRateExceeded(func() error {
		_, err = svc.CreateLogGroup(params)
		return err
	})

	return
}

func (tee *CWLogsTee) isStreamExist(svc cloudwatchlogsiface.CloudWatchLogsAPI) (exist bool, err error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(tee.LogGroupName),
		LogStreamNamePrefix: aws.String(tee.LogStreamName),
	}

	err = backoffRateExceeded(func() error {
		err = svc.DescribeLogStreamsPages(params, func(page *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
			for _, stream := range page.LogStreams {
				if *stream.LogStreamName == tee.LogStreamName {
					exist = true
					return false
				}
			}

			return !lastPage
		})

		return err
	})

	return
}

func (tee *CWLogsTee) createLogStream(svc cloudwatchlogsiface.CloudWatchLogsAPI) (err error) {
	params := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(tee.LogGroupName),
		LogStreamName: aws.String(tee.LogStreamName),
	}

	err = backoffRateExceeded(func() error {
		_, err = svc.CreateLogStream(params)
		return err
	})

	return
}

func (tee *CWLogsTee) putLogsEvents(svc cloudwatchlogsiface.CloudWatchLogsAPI, message string, sequenceToken *string) (nextToken *string, err error) {
	params := &cloudwatchlogs.PutLogEventsInput{
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(tee.Now().UnixNano() / int64(time.Millisecond)),
			},
		},
		LogGroupName:  aws.String(tee.LogGroupName),
		LogStreamName: aws.String(tee.LogStreamName),
	}

	if sequenceToken != nil {
		params.SequenceToken = sequenceToken
	}

	var resp *cloudwatchlogs.PutLogEventsOutput

	err = backoffRateExceeded(func() error {
		resp, err = svc.PutLogEvents(params)
		return err
	})

	if err == nil {
		nextToken = resp.NextSequenceToken
		return
	}

	matched, _ := regexp.MatchString(`^(DataAlreadyAcceptedException|InvalidSequenceTokenException):`, err.Error())

	if matched {
		re := regexp.MustCompile(`\bsequenceToken(?: is)?: (\S+)\b`)
		md := re.FindStringSubmatch(err.Error())

		if len(md) == 2 {
			nextToken, err = tee.putLogsEvents(svc, message, aws.String(md[1]))
		}
	}

	return
}

func (tee *CWLogsTee) put(svc cloudwatchlogsiface.CloudWatchLogsAPI, message string, sequenceToken *string) (nextToken *string, err error) {
	nextToken, err = tee.putLogsEvents(svc, message, sequenceToken)
	return
}

func (tee *CWLogsTee) Tee() (err error) {
	svc := cloudwatchlogs.New(session.New())
	var sequenceToken *string

	exist, err := tee.isGroupExist(svc)

	if err != nil {
		return
	}

	if !exist {
		err = tee.createLogGroup(svc)

		if err != nil {
			return
		}
	}

	exist, err = tee.isStreamExist(svc)

	if err != nil {
		return
	}

	if !exist {
		err = tee.createLogStream(svc)

		if err != nil {
			return
		}
	}

	err = tee.scan(func(line string) (err error) {
		sequenceToken, err = tee.put(svc, line, sequenceToken)

		if err != nil {
			return
		}

		fmt.Fprintln(tee.Out, line)

		return
	})

	return
}
