package service

import "time"

type Service struct {
	Name                string
	Urls                []string
	LbStrategy          string
	HealthCheckPath     string
	HealthCheckInterval time.Duration
}
