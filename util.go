package main

import (
	"flag"
)

// via https://github.com/skynetservices/skynet/blob/master/config.go
func getFlagName(f string) (name string) {
	if f[0] == '-' {
		minusCount := 1

		if f[1] == '-' {
			minusCount++
		}

		f = f[minusCount:]

		for i := 0; i < len(f); i++ {
			if f[i] == '=' || f[i] == ' ' {
				break
			}

			name += string(f[i])
		}
	}

	return
}

func splitFlagsetFromArgs(flagset *flag.FlagSet, args []string) (flagsetArgs []string, additionalArgs []string) {
	for _, f := range args {
		if flagset.Lookup(getFlagName(f)) != nil {
			flagsetArgs = append(flagsetArgs, f)
		} else {
			additionalArgs = append(additionalArgs, f)
		}
	}

	return
}
