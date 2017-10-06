package ui

import (
	"os"

	"github.com/cloudfoundry/cli/cf/i18n"
	"github.com/cloudfoundry/cli/cf/terminal"
)

var teePrinter *terminal.TeePrinter
var ui terminal.UI

func init() {
	i18n.T = func(translationID string, args ...interface{}) string {
		return translationID
	}
	teePrinter = terminal.NewTeePrinter()
	ui = terminal.NewUI(os.Stdin, teePrinter)
	teePrinter.DisableTerminalOutput(false)
}

func SetOutputBucket(bucket *[]string) {
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
	return ui.Ask(prompt, args...)
}

func Confirm(message string, args ...interface{}) bool {
	return ui.Confirm(message, args...)
}

func Ok() {
	ui.Ok()
}

func Failed(message string, args ...interface{}) {
	ui.Failed(message, args...)
}

func PanicQuietly() {
	ui.PanicQuietly()
}

func LoadingIndication() {
	ui.LoadingIndication()
}

func Table(headers []string) terminal.Table {
	return ui.Table(headers)
}
