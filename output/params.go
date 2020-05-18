package output

import (
	"github.com/spf13/cast"
	"time"
)

type Params map[string]interface{}

func (p Params) Add(key string, v interface{}) {
	p[key] = v
}

func (p Params) Get(key string) (v interface{}, ok bool) {
	v, ok = p[key]
	return
}

func (p Params) Exists(key string) bool {
	_, ok := p.Get(key)
	return ok
}

func (p Params) Bool(key string, def bool) bool {
	v, ok := p.Get(key)
	if !ok {
		return def
	}
	if x, err := cast.ToBoolE(v); err == nil {
		return x
	}
	return def
}

func (p Params) Int(key string, def int) int {
	v, ok := p.Get(key)
	if !ok {
		return def
	}
	if x, err := cast.ToIntE(v); err == nil {
		return x
	}
	return def
}

func (p Params) String(key string, def string) string {
	v, ok := p.Get(key)
	if !ok {
		return def
	}
	if x, err := cast.ToStringE(v); err == nil {
		return x
	}
	return def
}

func (p Params) Float(key string, def float64) float64 {
	v, ok := p.Get(key)
	if !ok {
		return def
	}
	if x, err := cast.ToFloat64E(v); err == nil {
		return x
	}
	return def
}

func (p Params) Duration(key string, def time.Duration) time.Duration {
	v, ok := p.Get(key)
	if !ok {
		return def
	}
	if x, err := cast.ToDurationE(v); err == nil {
		return x
	}
	return def
}