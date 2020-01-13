package testutil

import (
	"github.com/cloudfoundry-incubator/multiapps-cli-plugin/ui"
)

type OutputCapturer interface {
	CaptureOutputAndStatus(block func() int) ([]string, int)
	CaptureOutput(block func()) []string
}

type UIOutputCapturer struct{}

func NewUIOutputCapturer() OutputCapturer {
	return &UIOutputCapturer{}
}

func (oc *UIOutputCapturer) CaptureOutput(block func()) []string {
	var bucket []string
	ui.SetOutputBucket(&bucket)
	block()
	return bucket
}

func (oc *UIOutputCapturer) CaptureOutputAndStatus(block func() int) ([]string, int) {
	var bucket []string
	ui.SetOutputBucket(&bucket)
	status := block()
	return bucket, status
}
