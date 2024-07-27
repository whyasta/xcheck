package checks

type Check interface {
	Pass() bool
	Name() string
}

type FailureNotification struct {
	Threshold uint32
	Chan      chan error
}
