package driver

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrIncompatible = errors.New("driver is incompatible")
	ErrNotInterface = errors.New("driver is not an interface")
	ErrNotFound     = errors.New("driver not found")
)

type DriverProvider[T any] interface {
	Name() string
	CheckCompatibility(ctx context.Context) error
	New(ctx context.Context) (T, error)
}

var Drivers = map[reflect.Type]map[string]any{}

func Register[T any](provider DriverProvider[T]) {
	var t reflect.Type = reflect.TypeFor[T]()
	if t.Kind() != reflect.Interface {
		panic(fmt.Errorf("driver %s is not an interface", t.String()))
	}
	name := provider.Name()
	if _, ok := Drivers[t]; !ok {
		Drivers[t] = make(map[string]any)
	}
	if _, ok := Drivers[t][name]; ok {
		panic(fmt.Errorf("driver %s registered for %s being registered twice", name, t.String()))
	}
	Drivers[t][name] = provider
}

func GetByName[T any](name string) (DriverProvider[T], error) {
	var t reflect.Type = reflect.TypeFor[T]()
	if t.Kind() != reflect.Interface {
		return nil, ErrNotInterface
	}
	if providers, ok := Drivers[t]; ok {
		if provider, ok := providers[name]; ok {
			return provider.(DriverProvider[T]), nil
		}
	}
	return nil, ErrNotFound
}

func Get[T any](ctx context.Context) (T, error) {
	var t reflect.Type = reflect.TypeFor[T]()
	var zero T
	if t.Kind() != reflect.Interface {
		return zero, ErrNotInterface
	}
	if providers, ok := Drivers[t]; ok {
		if len(providers) == 0 {
			return zero, ErrNotFound
		}
		for _, provider := range providers {
			d := provider.(DriverProvider[T])
			return d.New(ctx)
		}
	}
	return zero, ErrNotFound
}
