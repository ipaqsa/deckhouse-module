package keyvalue

import "sync"

type Config struct {
	Enabled bool `mapstructure:"enabled"`
}
type Service struct {
	mtx  sync.Mutex
	data map[string]string
}

func New() *Service {
	return &Service{
		data: make(map[string]string),
	}
}

func (s *Service) Get(key string) string {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	return s.data[key]
}

func (s *Service) Set(key, value string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.data[key] = value
}

func (s *Service) Delete(key string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	delete(s.data, key)
}
