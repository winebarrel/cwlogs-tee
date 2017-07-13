package cwlogs_tee

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/bluele/go-timecop"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mockaws"
	"testing"
	"time"
)

func TestIsGroupExistWhenExist(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName: "my-group",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().DescribeLogGroupsPages(
		&cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String("my-group"),
		},
		gomock.Any(),
	).Do(func(_ *cloudwatchlogs.DescribeLogGroupsInput, fn func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool) {
		fn(
			&cloudwatchlogs.DescribeLogGroupsOutput{
				LogGroups: []*cloudwatchlogs.LogGroup{
					&cloudwatchlogs.LogGroup{
						LogGroupName: aws.String("my-group"),
					},
				},
			},
			true,
		)
	}).Return(
		nil,
	)

	exist, _ := tee.isGroupExist(mockcwlogs)

	assert.Equal(exist, true)
}

func TestIsGroupExistWhenNotExist(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName: "my-group",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().DescribeLogGroupsPages(
		&cloudwatchlogs.DescribeLogGroupsInput{
			LogGroupNamePrefix: aws.String("my-group"),
		},
		gomock.Any(),
	).Do(func(_ *cloudwatchlogs.DescribeLogGroupsInput, fn func(*cloudwatchlogs.DescribeLogGroupsOutput, bool) bool) {
		fn(
			&cloudwatchlogs.DescribeLogGroupsOutput{
				LogGroups: []*cloudwatchlogs.LogGroup{
					&cloudwatchlogs.LogGroup{
						LogGroupName: aws.String("my-group (untruth)"),
					},
				},
			},
			true,
		)
	}).Return(
		nil,
	)

	exist, _ := tee.isGroupExist(mockcwlogs)

	assert.Equal(exist, false)
}

func TestCreateLogGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName: "my-group",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().CreateLogGroup(
		&cloudwatchlogs.CreateLogGroupInput{
			LogGroupName: aws.String("my-group"),
		},
	).Return(
		nil,
		nil,
	)

	tee.createLogGroup(mockcwlogs)
}

func TestIsStreamExist(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().DescribeLogStreamsPages(
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String("my-group"),
			LogStreamNamePrefix: aws.String("my-stream"),
		},
		gomock.Any(),
	).Do(func(_ *cloudwatchlogs.DescribeLogStreamsInput, fn func(*cloudwatchlogs.DescribeLogStreamsOutput, bool) bool) {
		fn(
			&cloudwatchlogs.DescribeLogStreamsOutput{
				LogStreams: []*cloudwatchlogs.LogStream{
					&cloudwatchlogs.LogStream{
						LogStreamName: aws.String("my-stream"),
					},
				},
			},
			true,
		)
	}).Return(
		nil,
	)

	exist, _ := tee.isStreamExist(mockcwlogs)

	assert.Equal(exist, true)
}

func TestIsStreamNotExist(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().DescribeLogStreamsPages(
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String("my-group"),
			LogStreamNamePrefix: aws.String("my-stream"),
		},
		gomock.Any(),
	).Do(func(_ *cloudwatchlogs.DescribeLogStreamsInput, fn func(*cloudwatchlogs.DescribeLogStreamsOutput, bool) bool) {
		fn(
			&cloudwatchlogs.DescribeLogStreamsOutput{
				LogStreams: []*cloudwatchlogs.LogStream{
					&cloudwatchlogs.LogStream{
						LogStreamName: aws.String("my-stream (untruth)"),
					},
				},
			},
			true,
		)
	}).Return(
		nil,
	)

	exist, _ := tee.isStreamExist(mockcwlogs)

	assert.Equal(exist, false)
}

func TestCreateLogStream(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().CreateLogStream(
		&cloudwatchlogs.CreateLogStreamInput{
			LogGroupName:  aws.String("my-group"),
			LogStreamName: aws.String("my-stream"),
		},
	).Return(
		nil,
		nil,
	)

	tee.createLogStream(mockcwlogs)
}

func TestPutLogsEvents(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now, _ := time.Parse("2006-01-02 15:04:06", "2013-04-09 22:57:14")
	timecop.Freeze(now)

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
		Now:           timecop.Now,
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().PutLogEvents(
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Message:   aws.String("hello"),
					Timestamp: aws.Int64(1397084220000),
				},
			},
			LogGroupName:  aws.String("my-group"),
			LogStreamName: aws.String("my-stream"),
		},
	).Return(
		&cloudwatchlogs.PutLogEventsOutput{
			NextSequenceToken: aws.String("49559923757453780189052575064556876604280734704075802306"),
		},
		nil,
	)

	token, _ := tee.putLogsEvents(mockcwlogs, "hello", nil)

	assert.Equal(token, aws.String("49559923757453780189052575064556876604280734704075802306"))
}

func TestPutLogsEventsWithError(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now, _ := time.Parse("2006-01-02 15:04:06", "2013-04-09 22:57:14")
	timecop.Freeze(now)

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
		Now:           timecop.Now,
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().PutLogEvents(
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Message:   aws.String("hello"),
					Timestamp: aws.Int64(1397084220000),
				},
			},
			LogGroupName:  aws.String("my-group"),
			LogStreamName: aws.String("my-stream"),
		},
	).Return(
		nil,
		fmt.Errorf(`InvalidSequenceTokenException: The given sequenceToken is invalid. The next expected sequenceToken is: 49559923757453780189052575064556876604280734704075802306
status code: 400, request id: 5f6dd490-5983-11e6-839e-ab3896032135`),
	)

	mockcwlogs.EXPECT().PutLogEvents(
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Message:   aws.String("hello"),
					Timestamp: aws.Int64(1397084220000),
				},
			},
			LogGroupName:  aws.String("my-group"),
			LogStreamName: aws.String("my-stream"),
			SequenceToken: aws.String("49559923757453780189052575064556876604280734704075802306"),
		},
	).Return(
		&cloudwatchlogs.PutLogEventsOutput{
			NextSequenceToken: aws.String("49559923757453780189052575064556876604280734704075802307"),
		},
		nil,
	)

	token, _ := tee.putLogsEvents(mockcwlogs, "hello", nil)

	assert.Equal(token, aws.String("49559923757453780189052575064556876604280734704075802307"))
}

func TestPut(t *testing.T) {
	assert := assert.New(t)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now, _ := time.Parse("2006-01-02 15:04:06", "2013-04-09 22:57:14")
	timecop.Freeze(now)

	tee := &CWLogsTee{
		LogGroupName:  "my-group",
		LogStreamName: "my-stream",
		Now:           timecop.Now,
	}

	mockcwlogs := mockaws.NewMockCloudWatchLogsAPI(ctrl)

	mockcwlogs.EXPECT().PutLogEvents(
		&cloudwatchlogs.PutLogEventsInput{
			LogEvents: []*cloudwatchlogs.InputLogEvent{
				{
					Message:   aws.String("hello"),
					Timestamp: aws.Int64(1397084220000),
				},
			},
			LogGroupName:  aws.String("my-group"),
			LogStreamName: aws.String("my-stream"),
		},
	).Return(
		&cloudwatchlogs.PutLogEventsOutput{
			NextSequenceToken: aws.String("49559923757453780189052575064556876604280734704075802306"),
		},
		nil,
	)

	token, _ := tee.put(mockcwlogs, "hello", nil)

	assert.Equal(token, aws.String("49559923757453780189052575064556876604280734704075802306"))
}
