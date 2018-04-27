package health

import (
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	"github.com/InVisionApp/go-health/fakes"
	"github.com/InVisionApp/go-logger"
	"github.com/InVisionApp/go-logger/shims/testlog"
)

var (
	testCheckInterval = time.Duration(10) * time.Millisecond
)

// since we dont have before each in this testing framework...
func setupNewTestHealth() *Health {
	h := New()
	// replace with silent logger
	h.Logger = log.NewNoop()

	return h
}

func TestNew(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Should return a filled out instance of Health", func(t *testing.T) {
		h := setupNewTestHealth()

		Expect(h.configs).ToNot(BeNil())
		Expect(h.states).ToNot(BeNil())
		Expect(h.runners).ToNot(BeNil())
	})
}

func TestAddChecks(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path", func(t *testing.T) {
		h := setupNewTestHealth()
		testConfig := &Config{
			Name:     "foo",
			Checker:  &fakes.FakeICheckable{},
			Interval: testCheckInterval,
			Fatal:    false,
		}

		err := h.AddChecks([]*Config{testConfig})

		Expect(err).To(BeNil())
		Expect(h.configs).To(ContainElement(testConfig))
		Expect(len(h.configs)).To(Equal(1))
	})

	t.Run("Should error if healthcheck is already running", func(t *testing.T) {
		h := setupNewTestHealth()
		h.active.setTrue()
		err := h.AddChecks([]*Config{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ErrNoAddCfgWhenActive))
	})

	t.Run("Should error if passed in empty config slice", func(t *testing.T) {
		h := setupNewTestHealth()
		err := h.AddChecks([]*Config{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ErrEmptyConfigs))
	})
}

func TestAddCheck(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path", func(t *testing.T) {
		h := setupNewTestHealth()
		testConfig := &Config{
			Name:     "foo",
			Checker:  &fakes.FakeICheckable{},
			Interval: testCheckInterval,
			Fatal:    false,
		}

		err := h.AddCheck(testConfig)

		Expect(err).To(BeNil())
		Expect(h.configs).To(ContainElement(testConfig))
		Expect(len(h.configs)).To(Equal(1))
	})

	t.Run("Should error if healthcheck is already running", func(t *testing.T) {
		h := setupNewTestHealth()
		h.active.setTrue()
		err := h.AddCheck(&Config{})

		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(ErrNoAddCfgWhenActive))
	})
}

func TestDisableLogging(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Should set logger to noop logger", func(t *testing.T) {
		h := New()
		// Should initially be set to a basic logger
		Expect(h.Logger).To(BeEquivalentTo(log.NewSimple()))

		// Should set it to a noop logger
		h.DisableLogging()
		Expect(h.Logger).To(BeEquivalentTo(log.NewNoop()))
	})
}

func TestFailed(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Should return false if a fatally configured check hasn't errored", func(t *testing.T) {
		t.Run("Happy path", func(t *testing.T) {
			h := setupNewTestHealth()
			checker1 := &fakes.FakeICheckable{}
			checker1.StatusReturns(nil, nil)

			cfgs := []*Config{
				{
					Name:     "foo",
					Checker:  checker1,
					Interval: testCheckInterval,
					Fatal:    true,
				},
			}

			err := h.AddChecks(cfgs)
			Expect(err).ToNot(HaveOccurred())

			err = h.Start()
			Expect(err).ToNot(HaveOccurred())

			// More brittleness -- need to wait to ensure our checks have executed
			time.Sleep(time.Duration(15) * time.Millisecond)

			states, failed, err := h.State()
			Expect(err).ToNot(HaveOccurred())
			Expect(failed).To(BeFalse())
			Expect(states).To(HaveKey("foo"))

			Expect(h.Failed()).To(BeFalse())
		})

	})

	t.Run("Should return true if a fatally configured check has failed", func(t *testing.T) {
		h := setupNewTestHealth()
		checker1 := &fakes.FakeICheckable{}
		checker1.StatusReturns(nil, fmt.Errorf("things broke"))

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    true,
			},
		}

		err := h.AddChecks(cfgs)
		Expect(err).ToNot(HaveOccurred())

		err = h.Start()
		Expect(err).ToNot(HaveOccurred())

		// More brittleness -- need to wait to ensure our checks have executed
		time.Sleep(time.Duration(15) * time.Millisecond)

		states, failed, err := h.State()
		Expect(err).ToNot(HaveOccurred())
		Expect(failed).To(BeTrue())
		Expect(states).To(HaveKey("foo"))

		Expect(h.Failed()).To(BeTrue())

	})
}

func TestState(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path", func(t *testing.T) {
		h := setupNewTestHealth()
		checker1 := &fakes.FakeICheckable{}
		checker1.StatusReturns(nil, fmt.Errorf("things broke"))

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    false,
			},
		}

		err := h.AddChecks(cfgs)
		Expect(err).ToNot(HaveOccurred())

		err = h.Start()
		Expect(err).ToNot(HaveOccurred())

		// More brittleness -- need to wait to ensure our checks have executed
		time.Sleep(time.Duration(15) * time.Millisecond)

		states, failed, err := h.State()
		Expect(err).ToNot(HaveOccurred())
		Expect(failed).To(BeFalse())
		Expect(states).To(HaveKey("foo"))
		Expect(states["foo"].Err).To(Equal("things broke"))
	})

	t.Run("When a fatally-configured check fails and recovers, state should get updated accordingly", func(t *testing.T) {
		h := setupNewTestHealth()
		checker1 := &fakes.FakeICheckable{}
		checker1.StatusReturns(nil, fmt.Errorf("things broke"))

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    true,
			},
		}

		err := h.AddChecks(cfgs)
		Expect(err).ToNot(HaveOccurred())

		err = h.Start()
		Expect(err).ToNot(HaveOccurred())

		// More brittleness -- need to wait to ensure our checks have executed
		time.Sleep(time.Duration(15) * time.Millisecond)

		states, failed, err := h.State()
		Expect(err).ToNot(HaveOccurred())
		Expect(failed).To(BeTrue())
		Expect(states).To(HaveKey("foo"))
		Expect(states["foo"].Err).To(Equal("things broke"))

		// And now, let's let it recover
		checker1.StatusReturns(nil, nil)
		time.Sleep(time.Duration(15) * time.Millisecond)

		statesRecov, failedRecov, errRecov := h.State()
		Expect(errRecov).ToNot(HaveOccurred())
		Expect(failedRecov).To(BeFalse())
		Expect(statesRecov).To(HaveKey("foo"))
		Expect(statesRecov["foo"].Err).To(BeEmpty())

	})

}

func TestStart(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path", func(t *testing.T) {
		h := setupNewTestHealth()
		checker1 := &fakes.FakeICheckable{}
		checker2 := &fakes.FakeICheckable{}

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    false,
			},
			{
				Name:     "bar",
				Checker:  checker2,
				Interval: testCheckInterval,
				Fatal:    false,
			},
		}

		err := h.AddChecks(cfgs)
		Expect(err).ToNot(HaveOccurred())

		testLogger := testlog.NewTestLog()
		h.Logger = testLogger

		err = h.Start()
		Expect(err).ToNot(HaveOccurred())
		// Correct number of runners were created
		Expect(len(h.runners)).To(Equal(2))

		// Runners are created (and saved) based on their name
		for _, v := range cfgs {
			Expect(h.runners).To(HaveKey(v.Name))
		}

		// This is pretty brittle - will update if this is causing random test failures
		time.Sleep(time.Duration(15) * time.Millisecond)

		// Both runners should've ran
		Expect(checker1.StatusCallCount()).To(Equal(2), "Checker should have been executed")
		Expect(checker2.StatusCallCount()).To(Equal(2), "Checker should have been executed")

		// Both runners should've recorded their state
		Expect(h.states).To(HaveKey("foo"))
		Expect(h.states).To(HaveKey("bar"))

		// Ensure that logger was hit as expected
		Expect(testLogger.CallCount()).To(Equal(2))

		msgs := testLogger.Bytes()
		for _, cfg := range cfgs {
			Expect(msgs).To(ContainSubstring("Starting checker name=" + cfg.Name))
		}
	})

	t.Run("Should error if healthcheck already running", func(t *testing.T) {
		h := setupNewTestHealth()

		err := h.AddCheck(&Config{})
		Expect(err).ToNot(HaveOccurred())

		// Set the healthcheck state to active
		h.active.setTrue()
		err = h.Start()
		Expect(err).To(Equal(ErrAlreadyRunning))
	})
}

func TestStop(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path", func(t *testing.T) {
		testLogger := testlog.NewTestLog()
		h, cfgs, err := setupRunners(nil, testLogger)

		Expect(err).ToNot(HaveOccurred())
		Expect(h).ToNot(BeNil())

		// A bit brittle, but it'll do
		time.Sleep(time.Duration(15) * time.Millisecond)
		Expect(len(h.states)).To(Equal(2))

		err = h.Stop()
		Expect(err).ToNot(HaveOccurred())

		// Wait a bit to ensure goroutines have exited
		time.Sleep(15 * time.Millisecond)

		// Runners map should be reset
		Expect(h.runners).To(BeEmpty())

		// Ensure that logger captured the start and stop messages
		Expect(testLogger.CallCount()).To(Equal(6))

		for _, cfg := range cfgs {
			// 3rd and 4th message should indicate goroutine exit
			msgs := testLogger.Bytes()
			Expect(msgs).To(ContainSubstring("Stopping checker name=" + cfg.Name))

			Expect(msgs).To(ContainSubstring("Checker exiting name=" + cfg.Name))
		}

		// Expect state map to be reset
		Expect(len(h.states)).To(Equal(0))
	})

	t.Run("Should error if healthcheck is not running", func(t *testing.T) {
		h := setupNewTestHealth()

		err := h.AddCheck(&Config{})
		Expect(err).ToNot(HaveOccurred())

		// Set the healthcheck state to active
		h.active.setFalse()
		err = h.Stop()
		Expect(err).To(Equal(ErrAlreadyStopped))
	})
}

func TestStartRunner(t *testing.T) {
	RegisterTestingT(t)

	t.Run("Happy path - checkers do not fail", func(t *testing.T) {
		h, cfgs, err := setupRunners(nil, nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(h).ToNot(BeNil())

		// Brittle...
		time.Sleep(time.Duration(15) * time.Millisecond)

		// Did the ticker fire and create a state entry?
		for _, c := range cfgs {
			Expect(h.states).To(HaveKey(c.Name))
			Expect(h.states[c.Name].Status).To(Equal("ok"))
		}

		// Since nothing has failed, healthcheck should _not_ be in failed state
		Expect(h.failed.val()).To(BeFalse())
	})

	t.Run("Happy path - 1 checker fails (non-fatal)", func(t *testing.T) {
		checker1 := &fakes.FakeICheckable{}
		checker2 := &fakes.FakeICheckable{}
		checker2Error := errors.New("something failed")
		checker2.StatusReturns(nil, checker2Error)

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    false,
			},
			{
				Name:     "bar",
				Checker:  checker2,
				Interval: testCheckInterval,
				Fatal:    false,
			},
		}
		h, _, err := setupRunners(cfgs, nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(h).ToNot(BeNil())

		// Brittle...
		time.Sleep(time.Duration(15) * time.Millisecond)

		// Did the ticker fire and create a state entry?
		for _, c := range cfgs {
			Expect(h.states).To(HaveKey(c.Name))
		}

		// First checker should've succeeded
		Expect(h.states[cfgs[0].Name].Status).To(Equal("ok"))

		// Second checker should've failed
		Expect(h.states[cfgs[1].Name].Status).To(Equal("failed"))
		Expect(h.states[cfgs[1].Name].Err).To(Equal(checker2Error.Error()))

		// Since nothing has failed, healthcheck should _not_ be in failed state
		Expect(h.failed.val()).To(BeFalse())
	})

	t.Run("Happy path - 1 checker fails (fatal)", func(t *testing.T) {
		checker1 := &fakes.FakeICheckable{}
		checker2 := &fakes.FakeICheckable{}
		checker2Err := errors.New("something failed")
		checker2.StatusReturns(nil, checker2Err)

		cfgs := []*Config{
			{
				Name:     "foo",
				Checker:  checker1,
				Interval: testCheckInterval,
				Fatal:    false,
			},
			{
				Name:     "bar",
				Checker:  checker2,
				Interval: testCheckInterval,
				Fatal:    true,
			},
		}
		h, _, err := setupRunners(cfgs, nil)

		Expect(err).ToNot(HaveOccurred())
		Expect(h).ToNot(BeNil())

		// Brittle...
		time.Sleep(time.Duration(15) * time.Millisecond)

		// Did the ticker fire and create a state entry?
		for _, c := range cfgs {
			Expect(h.states).To(HaveKey(c.Name))
		}

		// First checker should've succeeded
		Expect(h.states[cfgs[0].Name].Status).To(Equal("ok"))

		// Second checker should've failed
		Expect(h.states[cfgs[1].Name].Status).To(Equal("failed"))
		Expect(h.states[cfgs[1].Name].Err).To(Equal(checker2Err.Error()))

		// Since second checker has failed fatally, global healthcheck state should be failed as well
		Expect(h.failed.val()).To(BeTrue())
	})
}
