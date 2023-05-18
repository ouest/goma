package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ouest/goma"
)

const version = "0.0.1"

func main() {
	historyCmd := flag.NewFlagSet("history", flag.ExitOnError)
	historyPage := historyCmd.Uint("p", 0, "history page")
	historyNum := historyCmd.Uint("n", 10, "number of history")

	lockCmd := flag.NewFlagSet("lock", flag.ExitOnError)
	lockAccount := lockCmd.String("a", "nobody", "Lock Account")
	toggleCmd := flag.NewFlagSet("toggle", flag.ExitOnError)
	toggleAccount := toggleCmd.String("a", "nobody", "Toggle Account")
	unlockCmd := flag.NewFlagSet("unlock", flag.ExitOnError)
	unlockAccount := unlockCmd.String("a", "nobody", "Unlock Account")

	usage := `
    goma : Control SESAME SmartLock by SESAME v2 API

    goma state
    goma history [-p HISTORY_PAGE] [-n NUMBER_OF_HISTORY]
    goma lock [-a ACCOUNT]
    goma toggle [-a ACCOUNT]
    goma unlock [-a ACCOUNT]

    HISTORY_PAGE: page number of history search results
    NUMBER_OF_HISTORY: number of history search results per page
    ACCOUNT: (lock|toggle|unlock) account
    `
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "state":
		var body = goma.State()
		fmt.Println(body)
	case "history":
		historyCmd.Parse(os.Args[2:])
		var body = goma.History(*historyPage, *historyNum)
		fmt.Println(body)
	case "lock":
		lockCmd.Parse(os.Args[2:])
		goma.Lock(*lockAccount)
	case "toggle":
		toggleCmd.Parse(os.Args[2:])
		goma.Toggle(*toggleAccount)
	case "unlock":
		unlockCmd.Parse(os.Args[2:])
		goma.Unlock(*unlockAccount)
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}
}
