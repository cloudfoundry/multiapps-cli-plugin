package fakes

import "github.com/SAP/cf-mta-plugin/util"

type fakeHttpGetExecutor struct{
  exchanges map[string]int
}

func NewFakeHttpGetExecutor(exchanges map[string]int) util.HttpSimpleGetExecutor{
  return &fakeHttpGetExecutor{exchanges: exchanges}
}

func (f fakeHttpGetExecutor) ExecuteGetRequest(url string) (int, error) {
  return f.exchanges[url], nil
}
