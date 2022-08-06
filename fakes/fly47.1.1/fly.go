package main

import (
	"fmt"
	"os"
)

func main() {
	if !argsEqualTo([]string{"--target", "mock", "--verbose", "status"}) {
		os.Exit(1)
	}
}

func argsEqualTo(expectedArgs []string) bool {
	actualArgs := os.Args[1:]

	if len(expectedArgs) != len(actualArgs) {
		printError(expectedArgs, actualArgs)
		fmt.Fprintf(os.Stderr, "Expected length is %v, actual length was %v.\n", len(expectedArgs), len(actualArgs))

		return false
	}

	for i, expected := range expectedArgs {
		actual := actualArgs[i]

		if expected != actual {
			printError(expectedArgs, actualArgs)
			fmt.Fprintf(os.Stderr, "Argument #%v was expected to be %v, but was actually %v.\n", i, expected, actual)

			return false
		}
	}

	return true
}

func printError(expectedArgs, actualArgs []string) {
	fmt.Fprintf(os.Stderr, "Expected arguments %v, but got %v.\n", expectedArgs, actualArgs)
}
