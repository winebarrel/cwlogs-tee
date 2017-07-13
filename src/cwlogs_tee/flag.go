package cwlogs_tee

import (
	"flag"
	"fmt"
	"os"
)

var version string

func ParseFlag(tee *CWLogsTee) (err error) {
	var showVersion bool

	flag.StringVar(&tee.LogGroupName, "g", "", "log group name")
	flag.StringVar(&tee.LogStreamName, "s", "", "log stream name")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.Parse()

	if showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

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
