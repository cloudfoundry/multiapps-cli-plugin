package fakes

import (
	"sync"

	"github.com/cloudfoundry/multiapps-cli-plugin/clients/models"
	"github.com/cloudfoundry/multiapps-cli-plugin/clients/mtaclient_v2"
)

type FakeMtaV2ClientOperations struct {
	GetMtasStub        func(name, namespace *string, spaceGuid string) ([]*models.Mta, error)
	getMtasMutex       sync.RWMutex
	getMtasArgsForCall []struct {
		name      *string
		namespace *string
		spaceGuid string
	}
	getMtasReturns struct {
		result1 []*models.Mta
		result2 error
	}
	getMtasReturnsOnCall map[int]struct {
		result1 []*models.Mta
		result2 error
	}

	GetMtasForThisSpaceStub        func(name, namespace *string) ([]*models.Mta, error)
	getMtasForThisSpaceMutex       sync.RWMutex
	getMtasForThisSpaceArgsForCall []struct {
		name      *string
		namespace *string
	}
	getMtasForThisSpaceReturns struct {
		result1 []*models.Mta
		result2 error
	}
	getMtasForThisSpaceReturnsOnCall map[int]struct {
		result1 []*models.Mta
		result2 error
	}

	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeMtaV2ClientOperations) GetMtas(name, namespace *string, spaceGuid string) ([]*models.Mta, error) {
	fake.getMtasMutex.Lock()
	ret, specificReturn := fake.getMtasReturnsOnCall[len(fake.getMtasArgsForCall)]
	fake.getMtasArgsForCall = append(fake.getMtasArgsForCall, struct {
		name      *string
		namespace *string
		spaceGuid string
	}{name, namespace, spaceGuid})
	fake.recordInvocation("GetMtas", []interface{}{name, namespace, spaceGuid})
	fake.getMtasMutex.Unlock()
	if fake.GetMtasStub != nil {
		return fake.GetMtasStub(name, namespace, spaceGuid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtasReturns.result1, fake.getMtasReturns.result2
}

func (fake *FakeMtaV2ClientOperations) GetMtasCallCount() int {
	fake.getMtasMutex.RLock()
	defer fake.getMtasMutex.RUnlock()
	return len(fake.getMtasArgsForCall)
}

func (fake *FakeMtaV2ClientOperations) GetMtasArgsForCall(i int) (*string, *string, string) {
	fake.getMtasMutex.RLock()
	defer fake.getMtasMutex.RUnlock()
	return fake.getMtasArgsForCall[i].name, fake.getMtasArgsForCall[i].namespace, fake.getMtasArgsForCall[i].spaceGuid
}

func (fake *FakeMtaV2ClientOperations) GetMtasReturns(result1 []*models.Mta, result2 error) {
	fake.GetMtasStub = nil
	fake.getMtasReturns = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaV2ClientOperations) GetMtasReturnsOnCall(i int, result1 []*models.Mta, result2 error) {
	fake.GetMtasStub = nil
	if fake.getMtasReturnsOnCall == nil {
		fake.getMtasReturnsOnCall = make(map[int]struct {
			result1 []*models.Mta
			result2 error
		})
	}
	fake.getMtasReturnsOnCall[i] = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaV2ClientOperations) GetMtasForThisSpace(name, namespace *string) ([]*models.Mta, error) {
	fake.getMtasForThisSpaceMutex.Lock()
	ret, specificReturn := fake.getMtasForThisSpaceReturnsOnCall[len(fake.getMtasForThisSpaceArgsForCall)]
	fake.getMtasForThisSpaceArgsForCall = append(fake.getMtasForThisSpaceArgsForCall, struct {
		name      *string
		namespace *string
	}{name, namespace})
	fake.recordInvocation("GetMtasForThisSpace", []interface{}{name, namespace})
	fake.getMtasForThisSpaceMutex.Unlock()
	if fake.GetMtasForThisSpaceStub != nil {
		return fake.GetMtasForThisSpaceStub(name, namespace)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getMtasForThisSpaceReturns.result1, fake.getMtasForThisSpaceReturns.result2
}

func (fake *FakeMtaV2ClientOperations) GetMtasForThisSpaceCallCount() int {
	fake.getMtasForThisSpaceMutex.RLock()
	defer fake.getMtasForThisSpaceMutex.RUnlock()
	return len(fake.getMtasForThisSpaceArgsForCall)
}

func (fake *FakeMtaV2ClientOperations) GetMtasForThisSpaceArgsForCall(i int) (*string, *string) {
	fake.getMtasForThisSpaceMutex.RLock()
	defer fake.getMtasForThisSpaceMutex.RUnlock()
	return fake.getMtasForThisSpaceArgsForCall[i].name, fake.getMtasForThisSpaceArgsForCall[i].namespace
}

func (fake *FakeMtaV2ClientOperations) GetMtasForThisSpaceReturns(result1 []*models.Mta, result2 error) {
	fake.GetMtasForThisSpaceStub = nil
	fake.getMtasForThisSpaceReturns = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaV2ClientOperations) GetMtasForThisSpaceReturnsOnCall(i int, result1 []*models.Mta, result2 error) {
	fake.GetMtasForThisSpaceStub = nil
	if fake.getMtasForThisSpaceReturnsOnCall == nil {
		fake.getMtasForThisSpaceReturnsOnCall = make(map[int]struct {
			result1 []*models.Mta
			result2 error
		})
	}
	fake.getMtasForThisSpaceReturnsOnCall[i] = struct {
		result1 []*models.Mta
		result2 error
	}{result1, result2}
}

func (fake *FakeMtaV2ClientOperations) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMtasMutex.RLock()
	defer fake.getMtasMutex.RUnlock()
	fake.getMtasForThisSpaceMutex.RLock()
	defer fake.getMtasForThisSpaceMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeMtaV2ClientOperations) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ mtaclient_v2.MtaV2ClientOperations = new(FakeMtaV2ClientOperations)
