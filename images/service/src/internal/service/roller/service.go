package roller

import (
	"math/rand"
	"strconv"
)

const defaultLimit = 10

type Config struct {
	Enabled bool `mapstructure:"enabled"`
	Limit   int  `mapstructure:"limit"`
}
type Service struct {
	limit int
}

func New(config Config) *Service {
	if config.Limit == 0 {
		config.Limit = defaultLimit
	}

	return &Service{
		limit: config.Limit,
	}
}

func (s *Service) RollDice() string {
	return strconv.Itoa(rand.Int() % s.limit)
}
