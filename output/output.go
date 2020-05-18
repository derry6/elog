package output

import (
    "fmt"
    "io"
    "sync"
)

var (
    constructorsMu sync.Mutex
    constructors   = map[string]Constructor{}
)

type Output interface {
    io.Closer
    Options() *Options
    Write(level Level, line []byte) error
    Sync() error
}

type Constructor func(opts ...Option) (Output, error)

func Register(name string, constructor Constructor) error {
    constructorsMu.Lock()
    if _, ok := constructors[name]; ok {
        constructorsMu.Unlock()
        return fmt.Errorf("output already exists")
    }
    constructors[name] = constructor
    constructorsMu.Unlock()
    return nil
}

func New(name string, opts ...Option) (Output, error) {
    constructorsMu.Lock()
    constructor, _ := constructors[name]
    constructorsMu.Unlock()
    if constructor != nil {
        return constructor(opts...)
    }
    return nil, fmt.Errorf("output not found")
}
