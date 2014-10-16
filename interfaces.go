package memcached_api

type netConnector interface {
	Read(p []byte) (n int, err error)
	Write(p []byte) (n int, err error)
	Close() error
}

type incrDecr func(delta int64) (int64, error)
type informer func() (map[string]uint, error)
