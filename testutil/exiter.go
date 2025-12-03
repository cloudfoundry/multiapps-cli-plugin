package testutil

import (
	"bytes"
	"strings"

	io_helpers "code.cloudfoundry.org/cli/v8/cf/util/testhelpers/io"
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

type OutputCapturer interface {
	CaptureOutputAndStatus(block func() int) ([]string, int)
	CaptureOutput(block func()) []string
}

type StdoutOutputCapturer struct{}

func NewStdoutOutputCapturer() OutputCapturer {
	return &StdoutOutputCapturer{}
}

func (oc *StdoutOutputCapturer) CaptureOutput(block func()) []string {
	return io_helpers.CaptureOutput(block)
}

func (oc *StdoutOutputCapturer) CaptureOutputAndStatus(block func() int) ([]string, int) {
	var status int
	output := io_helpers.CaptureOutput(func() {
		status = block()
	})
	return output, status
}

type UIOutputCapturer struct{}

func NewUIOutputCapturer() OutputCapturer {
	return &UIOutputCapturer{}
}

func (oc *UIOutputCapturer) CaptureOutput(block func()) []string {
	bucket := new(bytes.Buffer)
	ui.SetOutputBucket(bucket)
	block()
	return strings.Split(strings.TrimSpace(bucket.String()), "\n")
}

func (oc *UIOutputCapturer) CaptureOutputAndStatus(block func() int) ([]string, int) {
	var status int
	output := oc.CaptureOutput(func() {
		status = block()
	})
	return output, status
}
