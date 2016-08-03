package cwlogs_tee

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"mockaws"
	"testing"
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
