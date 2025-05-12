package server

import "sync/atomic"

type AtomicHealth struct {
	ready atomic.Bool
}

func NewAtomicHealth() *AtomicHealth {
	h := &AtomicHealth{}
	h.ready.Store(false)

	return h
}

func (h *AtomicHealth) SetReady(v bool) {
	h.ready.Store(v)
}

func (h *AtomicHealth) IsReady() bool {
	return h.ready.Load()
}
