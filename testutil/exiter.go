package testutil

import (
	io_helpers "github.com/cloudfoundry/cli/testhelpers/io"
	"github.com/SAP/cf-mta-plugin/ui"
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
	bucket := []string{}
	ui.SetOutputBucket(&bucket)
	block()
	return bucket
}

func (oc *UIOutputCapturer) CaptureOutputAndStatus(block func() int) ([]string, int) {
	bucket := []string{}
	ui.SetOutputBucket(&bucket)
	status := block()
	return bucket, status
}

var defaultOutputCapturer = NewStdoutOutputCapturer()
