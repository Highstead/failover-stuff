package failover

type lockService interface {
	AquireLock() bool
	Release() error
	HaveLock() bool
}
