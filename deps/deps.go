package deps

import (
	"fmt"
	"time"

	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-logger/shims/logrus"
	"github.com/InVisionApp/rye"
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/newrelic/go-agent"

	"github.com/talpert/helloaddon/config"
	"github.com/talpert/helloaddon/deps/backends"
)

type Dependencies struct {
	StatsD    statsd.Statter
	MWHandler *rye.MWHandler
	NRApp     newrelic.Application

	//DALs

	Backends *backends.Backends
	Health   health.IHealth
}

func New(cfg *config.Config) (*Dependencies, error) {
	gohealth := health.New()
	gohealth.Logger = logrus.New(nil)

	d := &Dependencies{
		Health: gohealth,
	}

	// StatsD
	if err := d.setupStatsdClient(cfg); err != nil {
		return nil, err
	}

	// Rye
	d.setupRyeMiddleware(cfg, d.StatsD)

	//Connect to backend DBs and APIs
	be, err := backends.NewBackends(cfg)
	if err != nil {
		return nil, err
	}

	d.Backends = be

	// Setup DALs and Managers
	// NOTE: All DALs must be created before Managers because Managers depend on DALs

	//dalHealthchecks, err := d.setupDALs(cfg)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if err := d.setupManagers(cfg); err != nil {
	//	return nil, err
	//}
	//
	// Health related calls should always be the last thing here
	//if err := d.Health.AddChecks(dalHealthchecks); err != nil {
	//	return nil, err
	//}
	//
	//if err := d.Health.AddChecks(d.Backends.Statuses); err != nil {
	//	return nil, err
	//}
	//
	//if err := d.Health.Start(); err != nil {
	//	return nil, err
	//}

	return d, nil
}

func (d *Dependencies) setupDALs(cfg *config.Config) ([]*health.Config, error) {
	return nil, nil
}

// Managers provide a way to abstract DAL's.
//
// _USE_ them if you have a bunch of DAL's that need to be called as part of the
// same transaction.
//
// _DO NOT_ use them only if you have just a couple of loosely related DAL's
// that can be used directly in the handler(s) (without creating a 200 line handler).
func (d *Dependencies) setupManagers(cfg *config.Config) error {
	// sm := smgr.New(d.SubsDAL, d.TrialDAL, d.ZuoraDAL)
	// d.SubsMgr = sm

	return nil
}

func (d *Dependencies) setupStatsdClient(cfg *config.Config) error {
	flushInterval := time.Duration(100 * time.Millisecond)

	statsdClient, err := statsd.NewBufferedClient(cfg.StatsDAddress, cfg.StatsDPrefix, flushInterval, 0)
	if err != nil {
		return fmt.Errorf("Unable to instantiate statsd client: %v", err)
	}

	d.StatsD = statsdClient

	return nil
}

func (d *Dependencies) setupRyeMiddleware(cfg *config.Config, statter statsd.Statter) {
	d.MWHandler = rye.NewMWHandler(rye.Config{
		Statter:  statter,
		StatRate: cfg.StatsDRate,
	})
}
