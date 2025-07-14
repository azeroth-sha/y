package ygrace

import (
	"time"

	"github.com/azeroth-sha/y/ylog"
)

type Service interface {
	Name() string
	Priority() int
	Serv(ylog.Logger) error
	Down(ylog.Logger) error
}

type ServWait interface {
	ServWait() time.Duration
}

type DownWait interface {
	DownWait() time.Duration
}
