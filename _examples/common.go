package examples

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// ExitCode represents an exit code
type ExitCode int

const (
	// ExitCodeSuccess used when process succeeded
	ExitCodeSuccess ExitCode = 0
	// ExitCodeWrongSyntax used when command called with wrong syntax
	ExitCodeWrongSyntax ExitCode = iota + 10
	// ExitCodeCouldNotOpenRepository used when the repository could not be opened
	ExitCodeCouldNotOpenRepository
	// ExitCodeCouldNotParseRevision used when the revision could not be parsed
	ExitCodeCouldNotParseRevision
	// ExitCodeCouldNotTraverseHistory used when the history could not be traversed
	ExitCodeCouldNotTraverseHistory
	// ExitCodeWrongCommitHash used when no commit found by hash
	ExitCodeWrongCommitHash
	// ExitCodeUnexpected used when process failed for unknown reason
	ExitCodeUnexpected
	// ExitCodeExpected used when process failed for known reason
	ExitCodeExpected
)

var errors = map[ExitCode]string{
	ExitCodeWrongSyntax:             "wrong syntax",
	ExitCodeCouldNotOpenRepository:  "not a git repository",
	ExitCodeCouldNotParseRevision:   "could not parse revision '%s'",
	ExitCodeCouldNotTraverseHistory: "could not traverse the repository history",
	ExitCodeWrongCommitHash:         "could not find commit '%s'",
	ExitCodeUnexpected:              "unhandled error",
	ExitCodeExpected:                "%s",
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		helpText := fmt.Sprintf("Usage: %s %s", "%_COMMAND_NAME_%", strings.Join(arg, " "))
		WrongSyntaxAndExit(helpText)
	}
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	ExitIfError(err, ExitCodeUnexpected)
}

// Print should be used to display a regular message
func Print(mainText string) {
	fmt.Printf("%s\n", mainText)
}

// Info should be used to describe the example commands that are about to run.
func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Warning should be used to display a warning
func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

// Error should be used to display an error
func Error(err error) {
	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
}

// HelpAndExit displays a help message and exits with ExitCodeSuccess
func HelpAndExit(desc, helpMsg string) {
	Print(desc)
	Print(strings.Replace(helpMsg, "%_COMMAND_NAME_%", os.Args[0], -1))
	os.Exit(int(ExitCodeSuccess))
}

// WrongSyntaxAndExit displays a wrong syntax message and exits with ExitCodeWrongSyntax
func WrongSyntaxAndExit(helpMsg string) {
	Error(fmt.Errorf(msg(ExitCodeWrongSyntax)))
	Print(strings.Replace(helpMsg, "%_COMMAND_NAME_%", os.Args[0], -1))
	os.Exit(int(ExitCodeWrongSyntax))
}

// ExitIfError stop the execution of the program with the passed exitCode, if
// passed an error; it will use the passed args to provide more info
func ExitIfError(err error, code ExitCode, args ...interface{}) {
	if err == nil {
		return
	}

	Error(fmt.Errorf("%s\n  %s", msg(code, args...), err))
	os.Exit(int(code))
}

func msg(code ExitCode, args ...interface{}) string {
	if txt, ok := errors[code]; ok {
		return fmt.Sprintf(txt, args...)
	}

	switch {
	case len(args) == 0:
		return errors[ExitCodeUnexpected]
	case len(args) == 1:
		return fmt.Sprintf("%s", args[0])
	}

	return fmt.Sprintf(fmt.Sprintf("%s", args[0]), args[1:]...)
}

// PrintCommits prints the commit hash
// if LOG_LEVEL env variable is set to `verbose` it will show also the commit message
func PrintCommits(commits ...*object.Commit) {
	for _, commit := range commits {
		printCommit(commit)
	}
}

func printCommit(commit *object.Commit) {
	if os.Getenv("LOG_LEVEL") == "verbose" {
		fmt.Printf(
			"\x1b[36;1m%s \x1b[90;21m%s\x1b[0m %s\n",
			commit.Hash.String()[:7],
			commit.Hash.String(),
			strings.Split(commit.Message, "\n")[0],
		)
	} else {
		Print(commit.Hash.String())
	}
}
