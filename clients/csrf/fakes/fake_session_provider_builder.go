package fakes

type FakeSessionProviderBuilder struct {
	fakeSessionProvider FakeSessionProvider
}

func NewFakeSessionProviderBuilder() *FakeSessionProviderBuilder {
	return &FakeSessionProviderBuilder{}
}

func (f *FakeSessionProviderBuilder) GetSession(err error) *FakeSessionProviderBuilder {
	f.fakeSessionProvider.GetSessionReturns(err)
	return f
}

func (f *FakeSessionProviderBuilder) Build() FakeSessionProvider {
	return f.fakeSessionProvider
}
