package sayer

const defaultText = "hello world!"

type Config struct {
	Enabled bool   `mapstructure:"enabled"`
	Text    string `mapstructure:"text"`
}
type Service struct {
	text string
}

func New(config Config) *Service {
	if config.Text == "" {
		config.Text = defaultText
	}

	return &Service{
		text: config.Text,
	}
}

func (s *Service) Say() string {
	return s.text
}
