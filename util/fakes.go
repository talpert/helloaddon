package util

//**Fake Mutex**

//go:generate counterfeiter -o ../fakes/syncfakes/fake_locker.go . Locker

type Locker interface {
	Lock()
	Unlock()
}
