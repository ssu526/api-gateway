package service

type Service struct {
	Name string   `yaml:"service"`
	Urls []string `yaml:"urls"`
}
