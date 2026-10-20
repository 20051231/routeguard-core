package hcore

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/sagernet/sing-box/adapter"
	"github.com/sagernet/sing-box/protocol/group"
)

type routeGuardPool struct {
	ctx       context.Context
	selector  *group.Selector
	terminals map[string]adapter.WarmTerminal
}

// Capture a particular running instance under the lifecycle lock. Never retain
// that lock while authenticating: network control re-enters Android policy code.
func captureRouteGuardPool(epoch int64) (*routeGuardPool, error) {
	static.lock.Lock()
	defer static.lock.Unlock()
	box := static.Box()
	if epoch <= 0 || box == nil {
		return nil, errors.New("RouteGuard core is unavailable")
	}
	owned, ok := box.Outbound().Outbound("routeguard-normal")
	selector, valid := owned.(*group.Selector)
	if !ok || !valid || selector.RouteGuardEpoch() != epoch {
		return nil, errors.New("RouteGuard core generation changed")
	}
	pool := &routeGuardPool{ctx: static.Context(), selector: selector, terminals: make(map[string]adapter.WarmTerminal)}
	for _, tag := range selector.All() {
		outbound, present := box.Outbound().Outbound(tag)
		if !present {
			return nil, errors.New("RouteGuard terminal is unavailable")
		}
		if terminal, warm := outbound.(adapter.WarmTerminal); warm {
			pool.terminals[tag] = terminal
		}
	}
	if pool.ctx == nil || pool.ctx.Err() != nil {
		return nil, errors.New("RouteGuard core generation ended")
	}
	return pool, nil
}

type routeGuardTerminalState struct {
	Tag string `json:"tag"`
	adapter.WarmSession
}
type routeGuardPoolState struct {
	Epoch     int64                     `json:"epoch"`
	Selected  string                    `json:"selected"`
	Terminals []routeGuardTerminalState `json:"terminals"`
}

func (p *routeGuardPool) snapshot(epoch int64) (string, error) {
	if p.ctx.Err() != nil {
		return "", errors.New("RouteGuard core generation ended")
	}
	state := routeGuardPoolState{Epoch: epoch, Selected: p.selector.Selected("tcp").Tag(), Terminals: make([]routeGuardTerminalState, 0, len(p.terminals))}
	// Preserve the trusted configured order; never export profile endpoints/auth.
	for _, tag := range p.selector.All() {
		if terminal, ok := p.terminals[tag]; ok {
			state.Terminals = append(state.Terminals, routeGuardTerminalState{Tag: tag, WarmSession: terminal.WarmSession()})
		}
	}
	encoded, err := json.Marshal(state)
	if p.ctx.Err() != nil {
		return "", errors.New("RouteGuard core generation ended")
	}
	return string(encoded), err
}

// These methods are exposed only through the in-process Mobile bridge, not
// through unauthenticated loopback gRPC/Clash mutation endpoints.
func RouteGuardPoolSnapshot(epoch int64) (string, error) {
	pool, err := captureRouteGuardPool(epoch)
	if err != nil {
		return "", err
	}
	return pool.snapshot(epoch)
}

func RouteGuardWarmPool(epoch int64) (string, error) {
	pool, err := captureRouteGuardPool(epoch)
	if err != nil {
		return "", err
	}
	ctx, cancel := context.WithTimeout(pool.ctx, 8*time.Second)
	defer cancel()
	var pending sync.WaitGroup
	for _, terminal := range pool.terminals {
		pending.Add(1)
		go func(terminal adapter.WarmTerminal) { defer pending.Done(); _ = terminal.Warm(ctx) }(terminal)
	}
	pending.Wait()
	return pool.snapshot(epoch)
}

func RouteGuardSelectKeep(epoch int64, tag string) error {
	static.lock.Lock()
	defer static.lock.Unlock()
	box := static.Box()
	if box == nil {
		return errors.New("RouteGuard core is unavailable")
	}
	owned, ok := box.Outbound().Outbound("routeguard-normal")
	selector, valid := owned.(*group.Selector)
	if !ok || !valid || !selector.RouteGuardSelectKeep(epoch, tag) {
		return errors.New("RouteGuard selection or generation is not approved")
	}
	return nil
}
