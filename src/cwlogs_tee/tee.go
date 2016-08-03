package cwlogs_tee

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"io"
	"regexp"
	"strings"
	"time"
)

type CWLogsTee struct {
	LogGroupName  string
	LogStreamName string
	In            io.Reader
	Out           io.Writer
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

func (tee *CWLogsTee) isGroupExist(svc *cloudwatchlogs.CloudWatchLogs) (exist bool, err error) {

	params := &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(tee.LogGroupName),
		Limit:              aws.Int64(1),
	}

	resp, err := svc.DescribeLogGroups(params)

	if err != nil {
		return
	}

	exist = len(resp.LogGroups) > 0

	return
}

func (tee *CWLogsTee) createLogGroup(svc *cloudwatchlogs.CloudWatchLogs) (err error) {
	params := &cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(tee.LogGroupName),
	}

	_, err = svc.CreateLogGroup(params)

	return
}

func (tee *CWLogsTee) isStreamExist(svc *cloudwatchlogs.CloudWatchLogs) (exist bool, err error) {
	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(tee.LogGroupName),
		LogStreamNamePrefix: aws.String(tee.LogStreamName),
		Limit:               aws.Int64(1),
	}

	err = svc.DescribeLogStreamsPages(params, func(page *cloudwatchlogs.DescribeLogStreamsOutput, lastPage bool) bool {
		for _, stream := range page.LogStreams {
			if *stream.LogStreamName == tee.LogStreamName {
				exist = true
				return false
			}
		}

		return !lastPage
	})

	if err != nil {
		return
	}

	return
}

func (tee *CWLogsTee) createLogStream(svc *cloudwatchlogs.CloudWatchLogs) (err error) {
	params := &cloudwatchlogs.CreateLogStreamInput{
		LogGroupName:  aws.String(tee.LogGroupName),
		LogStreamName: aws.String(tee.LogStreamName),
	}

	_, err = svc.CreateLogStream(params)

	return
}

func (tee *CWLogsTee) putLogsEvents(svc *cloudwatchlogs.CloudWatchLogs, message string, sequenceToken *string) (nextToken *string, err error) {
	params := &cloudwatchlogs.PutLogEventsInput{
		LogEvents: []*cloudwatchlogs.InputLogEvent{
			{
				Message:   aws.String(message),
				Timestamp: aws.Int64(time.Now().UnixNano() / int64(time.Millisecond)),
			},
		},
		LogGroupName:  aws.String(tee.LogGroupName),
		LogStreamName: aws.String(tee.LogStreamName),
	}

	if sequenceToken != nil {
		params.SequenceToken = sequenceToken
	}

	resp, err := svc.PutLogEvents(params)

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

func (tee *CWLogsTee) put(svc *cloudwatchlogs.CloudWatchLogs, message string, sequenceToken *string) (nextToken *string, err error) {
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

	nextToken, err = tee.putLogsEvents(svc, message, sequenceToken)

	return
}

func (tee *CWLogsTee) Tee() (err error) {
	svc := cloudwatchlogs.New(session.New())
	var sequenceToken *string

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
