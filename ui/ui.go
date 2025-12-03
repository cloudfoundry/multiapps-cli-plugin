package ui

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"code.cloudfoundry.org/cli/v8/cf/i18n"
	"code.cloudfoundry.org/cli/v8/cf/terminal"
	"code.cloudfoundry.org/cli/v8/cf/trace"
)

var teePrinter *terminal.TeePrinter
var ui terminal.UI

func init() {
	i18n.T = func(translationID string, args ...interface{}) string {
		return translationID
	}
	disableColorsIfNeeded()
	teePrinter = terminal.NewTeePrinter(os.Stdout)
	ui = terminal.NewUI(os.Stdin, os.Stdout, teePrinter, trace.NewWriterPrinter(io.Discard, false))
}

func disableColorsIfNeeded() {
	if runtime.GOOS == "windows" {
		terminal.UserAskedForColors = "false"
		terminal.InitColorSupport()
	}
}

func SetOutputBucket(bucket io.Writer) {
	teePrinter.SetOutputBucket(bucket)
}

func ClearOutputBucket() {
	teePrinter.SetOutputBucket(nil)
}

func DisableTerminalOutput(disable bool) {
	teePrinter.DisableTerminalOutput(disable)
}

func PrintPaginator(rows []string, err error) {
	ui.PrintPaginator(rows, err)
}

func Say(message string, args ...interface{}) {
	ui.Say(message, args...)
}

func PrintCapturingNoOutput(message string, args ...interface{}) {
	ui.PrintCapturingNoOutput(message, args...)
}

func Warn(message string, args ...interface{}) {
	ui.Warn(message, args...)
}

func Ask(prompt string, args ...interface{}) (answer string) {
	return ui.Ask(fmt.Sprintf(prompt, args...))
}

func Confirm(message string, args ...interface{}) bool {
	return ui.Confirm(fmt.Sprintf(message, args...))
}

func Ok() {
	ui.Ok()
}

func Failed(message string, args ...interface{}) {
	ui.Failed(message, args...)
}

func LoadingIndication() {
	ui.LoadingIndication()
}

func Table(headers []string) *terminal.UITable {
	return ui.Table(headers)
}
