// Code generated by counterfeiter. DO NOT EDIT.
package depsfakes

import (
	"sync"

	"github.com/talpert/helloaddon/api/deps"
)

type FakeICheckable struct {
	CheckStatusStub        func() (interface{}, error)
	checkStatusMutex       sync.RWMutex
	checkStatusArgsForCall []struct{}
	checkStatusReturns     struct {
		result1 interface{}
		result2 error
	}
	checkStatusReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	IsFatalStub        func() bool
	isFatalMutex       sync.RWMutex
	isFatalArgsForCall []struct{}
	isFatalReturns     struct {
		result1 bool
	}
	isFatalReturnsOnCall map[int]struct {
		result1 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeICheckable) CheckStatus() (interface{}, error) {
	fake.checkStatusMutex.Lock()
	ret, specificReturn := fake.checkStatusReturnsOnCall[len(fake.checkStatusArgsForCall)]
	fake.checkStatusArgsForCall = append(fake.checkStatusArgsForCall, struct{}{})
	fake.recordInvocation("CheckStatus", []interface{}{})
	fake.checkStatusMutex.Unlock()
	if fake.CheckStatusStub != nil {
		return fake.CheckStatusStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.checkStatusReturns.result1, fake.checkStatusReturns.result2
}

func (fake *FakeICheckable) CheckStatusCallCount() int {
	fake.checkStatusMutex.RLock()
	defer fake.checkStatusMutex.RUnlock()
	return len(fake.checkStatusArgsForCall)
}

func (fake *FakeICheckable) CheckStatusReturns(result1 interface{}, result2 error) {
	fake.CheckStatusStub = nil
	fake.checkStatusReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeICheckable) CheckStatusReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.CheckStatusStub = nil
	if fake.checkStatusReturnsOnCall == nil {
		fake.checkStatusReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.checkStatusReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeICheckable) IsFatal() bool {
	fake.isFatalMutex.Lock()
	ret, specificReturn := fake.isFatalReturnsOnCall[len(fake.isFatalArgsForCall)]
	fake.isFatalArgsForCall = append(fake.isFatalArgsForCall, struct{}{})
	fake.recordInvocation("IsFatal", []interface{}{})
	fake.isFatalMutex.Unlock()
	if fake.IsFatalStub != nil {
		return fake.IsFatalStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isFatalReturns.result1
}

func (fake *FakeICheckable) IsFatalCallCount() int {
	fake.isFatalMutex.RLock()
	defer fake.isFatalMutex.RUnlock()
	return len(fake.isFatalArgsForCall)
}

func (fake *FakeICheckable) IsFatalReturns(result1 bool) {
	fake.IsFatalStub = nil
	fake.isFatalReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeICheckable) IsFatalReturnsOnCall(i int, result1 bool) {
	fake.IsFatalStub = nil
	if fake.isFatalReturnsOnCall == nil {
		fake.isFatalReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isFatalReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *FakeICheckable) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkStatusMutex.RLock()
	defer fake.checkStatusMutex.RUnlock()
	fake.isFatalMutex.RLock()
	defer fake.isFatalMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeICheckable) recordInvocation(key string, args []interface{}) {
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

var _ deps.ICheckable = new(FakeICheckable)
