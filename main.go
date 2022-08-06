package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/concourse/fly/rc"
	"github.com/jessevdk/go-flags"
	"github.com/suhlig/apron-bus/atc"
)

type Args struct {
	TargetName rc.TargetName `short:"t" long:"target" description:"Concourse target name"`
	Verbose    bool          `long:"verbose" description:"Print API requests and responses"`
	URL        string        `short:"u" long:"url" description:"Concourse target URL"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("No args given.")
		fmt.Println("")
		fmt.Println("Usage: apron-bus selects the flyX.Y.Z binary that matches the target server's version; passing all arguments to it.")
		os.Exit(0)
	}

	flyTargets, err := rc.LoadTargets()

	if err != nil {
		fmt.Fprintf(os.Stderr, "apron-bus: could load targets: %v\n", err)
		os.Exit(1)
	}

	var apiURL string
	apronBusArgs := parseArgs()

	if apronBusArgs.TargetName != "" {
		if apronBusArgs.Verbose {
			fmt.Fprintf(os.Stderr, "apron-bus: looking up API endpoint for target %v\n", apronBusArgs.TargetName)
		}

		target, ok := flyTargets.Targets[apronBusArgs.TargetName]

		if !ok {
			fmt.Fprintf(os.Stderr, "apron-bus: could not find target %v\n", apronBusArgs)
			os.Exit(1)
		}

		apiURL = target.API
	} else if apronBusArgs.URL != "" {
		apiURL = apronBusArgs.URL
	} else {
		fmt.Fprintln(os.Stderr, "apron-bus: neither target name nor URL given; please use one of the fly commands directly.")
		os.Exit(1)
	}

	if apronBusArgs.Verbose {
		fmt.Fprintf(os.Stderr, "apron-bus: retrieving server version from API endpoint %v\n", apiURL)
	}

	client := atc.NewVersionClient(apiURL)
	version, err := client.GetServerVersion()

	if err != nil {
		fmt.Fprintf(os.Stderr, "apron-bus: %v\n", err)
		os.Exit(1)
	}

	if apronBusArgs.Verbose {
		fmt.Fprintf(os.Stderr, "apron-bus: server version is %v\n", version)
	}

	flyWithVersion := fmt.Sprintf("fly%v", version)
	pathToFly, err := exec.LookPath(flyWithVersion)

	if err != nil {
		fmt.Fprintf(os.Stderr, "apron-bus: could not find %v in $PATH (%v)\n", flyWithVersion, os.Getenv("PATH"))
		os.Exit(1)
	}

	args := os.Args
	args[0] = flyWithVersion

	if apronBusArgs.Verbose {
		fmt.Fprintf(os.Stderr, "apron-bus: invoking %v\n", strings.Join(args, " "))
	}

	if err := syscall.Exec(pathToFly, args, os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "apron-bus: could not invoke %v: %v\n", pathToFly, err)
		os.Exit(1)
	}
}

func parseArgs() *Args {
	args := Args{}

	parser := flags.NewParser(&args, flags.PassDoubleDash)
	parser.NamespaceDelimiter = "-"

	argsCopy := os.Args
	parser.ParseArgs(argsCopy) // ignoring errors because fly will handle them

	return &args
}
