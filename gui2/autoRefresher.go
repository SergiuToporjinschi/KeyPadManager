package gui

import (
	"context"
	"main/logger"
	"time"
)

type AutoRefresher struct {
	cancelCtx   context.CancelFunc
	running     bool
	runFunc     func()
	afterToggle func(running bool)
}

func NewAutoRefresher(runFunc func(), afterToggle func(bool)) *AutoRefresher {
	return &AutoRefresher{
		runFunc:     runFunc,
		afterToggle: afterToggle,
	}
}

func (ar *AutoRefresher) Start() {
	if ar.running {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	ar.cancelCtx = cancel
	ar.running = true
	go ar.autoRefresh(ctx)
}

func (ar *AutoRefresher) Stop() {
	if !ar.running {
		return
	}

	ar.cancelCtx()
	ar.running = false
}

func (ar *AutoRefresher) Toggle() {
	if ar.running {
		ar.Stop()
	} else {
		ar.Start()
	}
	if ar.afterToggle != nil {
		ar.afterToggle(ar.running)
	}
}

func (ar *AutoRefresher) autoRefresh(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Log.Debug("Stopping auto refresh")
			return
		case t := <-ticker.C:
			logger.Log.Debugf("Last refresh at %v", t)
			go ar.runFunc()
		}
	}
}
