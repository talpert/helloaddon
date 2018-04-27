package backends

import (
	"github.com/sirupsen/logrus"

	"github.com/InVisionApp/go-health"
	"github.com/talpert/helloaddon/config"
)

var (
	log *logrus.Entry
)

func init() {
	log = logrus.WithField("pkg", "backends")
}

type Backends struct {
	//Direct connection backends

	//API clients

	//backends dependency health
	Statuses []*health.Config

	//Use this to determine if the connect method has been
	// run successfully and it is safe to use the underlying backends.
	// Not the best way to do this but good enough for now.
	connected bool
}

func NewBackends(cfg *config.Config) (*Backends, error) {
	b := &Backends{
		Statuses:  []*health.Config{},
		connected: false, //ensure
	}

	//Connect to DBs

	// Clients
	b.connected = true

	return b, nil
}

func (b *Backends) IsConnected() bool {
	return b.connected
}
