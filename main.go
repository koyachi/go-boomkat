package main

import (
	"flag"
	"fmt"
)

var (
	version     = "0.0.1"
	showVersion = flag.Bool("version", false, "print version string")
)

func main() {
	flag.Parse()

	if *showVersion {
		fmt.Printf("boomkat %s\n", version)
		return
	}

	switch flag.Arg(0) {
	case "search":
		search(flag.Arg(1))
		return
	case "download":
		downloadTrack(flag.Arg(1), flag.Arg(2))
		return
	}
}
