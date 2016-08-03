package cwlogs_tee

import (
	"flag"
	"fmt"
)

func ParseFlag(tee *CWLogsTee) (err error) {
	flag.StringVar(&tee.LogGroupName, "g", "", "log group name")
	flag.StringVar(&tee.LogStreamName, "s", "", "log stream name")
	flag.Parse()

	if tee.LogGroupName == "" {
		err = fmt.Errorf("'-g' is required")
		return
	}

	if tee.LogStreamName == "" {
		err = fmt.Errorf("'-s' is required")
		return
	}

	return
}
