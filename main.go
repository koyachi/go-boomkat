package main

import (
	"flag"
	"fmt"
	"github.com/koyachi/go-boomkat/boomkat"
	"os"
)

var (
	version     = "0.0.1"
	flagset     = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	showVersion = flagset.Bool("version", false, "print version string")
	directory   = flagset.String("directory", "", "download directory")
	recordId    = flagset.String("record_id", "", "download record id")
	trackId     = flagset.String("track_id", "", "download track id")
)

func main() {
	flagsetArgs, additionalArgs := splitFlagsetFromArgs(flagset, os.Args[1:])

	flagset.Parse(flagsetArgs)

	if *showVersion {
		fmt.Printf("boomkat %s\n", version)
		return
	}

	if *directory != "" {
		boomkat.SetBoomkatDir(*directory)
	}
	fmt.Printf("boomkatDir =  %s\n", boomkat.BoomkatDir())
	//fmt.Printf("additionalArgs =  %v\n", additionalArgs)

	switch additionalArgs[0] {
	case "search":
		search(additionalArgs[1])
		return
	case "download":
		if *recordId == "" && len(additionalArgs) >= 2 {
			*recordId = additionalArgs[1]
			if *trackId == "" && len(additionalArgs) >= 3 {
				*trackId = additionalArgs[2]
			}
		} else {
			if *trackId == "" && len(additionalArgs) >= 2 {
				*trackId = additionalArgs[1]
			}
		}
		if *recordId != "" {
			if *trackId != "" {
				downloadTrack(*recordId, *trackId)
			} else {
				downloadRecord(*recordId)
			}
		} else {
			// TODO: display error message
		}
		return
	}
}
