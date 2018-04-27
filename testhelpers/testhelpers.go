package testhelpers

import (
	"github.com/newrelic/go-agent"
)

func NewDisabledTestGoAgent() newrelic.Application {
	conf := newrelic.NewConfig("helloaddon-tests", "")
	conf.Enabled = false

	fakeAgent, agentErr := newrelic.NewApplication(conf)
	if agentErr != nil {
		panic(agentErr)
	}

	return fakeAgent
}
