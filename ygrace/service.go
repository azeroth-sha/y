package ygrace

import (
	"time"

	"github.com/azeroth-sha/y/logger"
)

type Service interface {
	Name() string
	Priority() int
	Serv(logger.Logger) error
	Down(logger.Logger) error
}

type ServWait interface {
	ServWait() time.Duration
}

type DownWait interface {
	DownWait() time.Duration
}
