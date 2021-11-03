package testutil

import (
	. "github.com/onsi/gomega"
)

type Expector interface {
	ExpectNoError(err error)
	ExpectNoErrorAndResult(err error, result, expectedResult interface{})
	ExpectError(err error)
	ExpectErrorAndZeroResult(err error, result interface{})
	ExpectSuccess(status int, output []string)
	ExpectSuccessWithOutput(status int, output, expectedOutput []string)
	ExpectSuccessAndResult(status int, output []string, result, expectedResult interface{})
	ExpectFailure(status int, output []string, message string)
	ExpectFailureOnLine(status int, output []string, message string, line int)
	ExpectNonZeroStatus(status int)
	ExpectMessageOnLine(output []string, message string, line int)
}

type BaseExpector struct{}

func (ex *BaseExpector) ExpectNoError(err error) {
	Expect(err).ToNot(HaveOccurred())
}

func (ex *BaseExpector) ExpectNoErrorAndResult(err error, result, expectedResult interface{}) {
	ex.ExpectNoError(err)
	Expect(result).To(Equal(expectedResult))
}

func (ex *BaseExpector) ExpectError(err error) {
	Expect(err).To(HaveOccurred())
}

func (ex *BaseExpector) ExpectErrorAndZeroResult(err error, result interface{}) {
	Expect(err).To(HaveOccurred())
	Expect(result).To(BeZero())
}

type UIExpector struct{ BaseExpector }

func NewUIExpector() Expector {
	return &UIExpector{}
}

func (ex *UIExpector) ExpectSuccess(status int, output []string) {
	ex.ExpectSuccessWithOutput(status, output, []string{})
}

func (ex *UIExpector) ExpectSuccessWithOutput(status int, output, expectedOutput []string) {
	Expect(status).To(BeZero())
	Expect(output).To(Equal(expectedOutput))
}

func (ex *UIExpector) ExpectSuccessAndResult(status int, output []string, result, expectedResult interface{}) {
	ex.ExpectSuccess(status, output)
	Expect(result).To(Equal(expectedResult))
}

func (ex *UIExpector) ExpectFailure(status int, output []string, message string) {
	ex.ExpectFailureOnLine(status, output, message, 1)
}

func (ex *UIExpector) ExpectFailureOnLine(status int, output []string, message string, line int) {
	ex.ExpectNonZeroStatus(status)
	Expect(output[line]).To(ContainSubstring(message))
}

func (ex *UIExpector) ExpectNonZeroStatus(status int) {
	Expect(status).NotTo(BeZero())
}

func (ex *UIExpector) ExpectMessageOnLine(output []string, message string, line int) {
	Expect(output[line]).To(ContainSubstring(message))
}

var defaultExpector = NewUIExpector()

func ExpectNoError(err error) {
	defaultExpector.ExpectNoError(err)
}

func ExpectNoErrorAndResult(err error, result, expectedResult interface{}) {
	defaultExpector.ExpectNoErrorAndResult(err, result, expectedResult)
}

func ExpectError(err error) {
	defaultExpector.ExpectError(err)
}

func ExpectErrorAndZeroResult(err error, result interface{}) {
	defaultExpector.ExpectErrorAndZeroResult(err, result)
}

func ExpectSuccess(status int, output []string) {
	defaultExpector.ExpectSuccess(status, output)
}

func ExpectSuccessWithOutput(status int, output, expectedOutput []string) {
	defaultExpector.ExpectSuccessWithOutput(status, output, expectedOutput)
}

func ExpectSuccessAndResult(status int, output []string, result, expectedResult interface{}) {
	defaultExpector.ExpectSuccessAndResult(status, output, result, expectedResult)
}

func ExpectFailure(status int, output []string, message string) {
	defaultExpector.ExpectFailure(status, output, message)
}

func ExpectFailureOnLine(status int, output []string, message string, line int) {
	defaultExpector.ExpectFailureOnLine(status, output, message, line)
}
