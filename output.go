package elog

import (
	"fmt"
	"io"
	"sync"
)

var (
	_outsMu sync.Mutex
	_outs   = map[string]Constructor{}
)

type Output interface {
	io.Closer
	Write(level Level, line []byte) error
	Sync() error
}

type Constructor func(cfg *OutputConfig) (Output, error)

func RegisterOutput(name string, constructor Constructor) error {
	_outsMu.Lock()
	if _, ok := _outs[name]; ok {
		_outsMu.Unlock()
		return fmt.Errorf("output already exists")
	}
	_outs[name] = constructor
	_outsMu.Unlock()
	return nil
}

func newOutput(name string, cfg *OutputConfig) (Output, error) {
	_outsMu.Lock()
	constructor, _ := _outs[name]
	_outsMu.Unlock()
	if constructor != nil {
		return constructor(cfg)
	}
	return nil, fmt.Errorf("output not found")
}
