package main

import (
	"cwlogs_tee"
	"log"
	"os"
	"time"
)

func init() {
	log.SetFlags(0)
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	tee := &cwlogs_tee.CWLogsTee{
		In:  os.Stdin,
		Out: os.Stdout,
		Now: time.Now,
	}

	err := cwlogs_tee.ParseFlag(tee)

	if err != nil {
		panic(err)
	}

	err = tee.Tee()

	if err != nil {
		panic(err)
	}
}
