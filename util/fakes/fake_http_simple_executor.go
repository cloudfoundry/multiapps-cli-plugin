package fakes

import "github.com/SAP/cf-mta-plugin/util"

type fakeHttpGetExecutor struct{
  statusCode int
  err error
}

func NewFakeHttpGetExecutor(statusCode int, err error) util.HttpSimpleGetExecutor{
  return &fakeHttpGetExecutor{statusCode: statusCode, err: err}
}

func (f fakeHttpGetExecutor) ExecuteGetRequest(url string) (int, error) {
  return f.statusCode, f.err
}
