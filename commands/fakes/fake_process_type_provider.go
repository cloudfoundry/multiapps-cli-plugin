package fakes

type FakeDeployCommandProcessTypeProvider struct {}

func (FakeDeployCommandProcessTypeProvider) GetProcessType() string {
	return "DEPLOY"
}

type FakeBlueGreenCommandProcessTypeProvider struct {}

func (FakeBlueGreenCommandProcessTypeProvider) GetProcessType() string {
	return "BLUE_GREEN_DEPLOY"
}
