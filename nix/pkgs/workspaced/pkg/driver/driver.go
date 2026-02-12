package driver

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
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

type doctorEntry struct {
	InterfaceName string
	DriverName    string
	Check         func(context.Context) error
}

var (
	mu         sync.RWMutex
	Drivers    = map[reflect.Type]map[string]any{}
	doctorList = []doctorEntry{}
)

func Register[T any](provider DriverProvider[T]) {
	mu.Lock()
	defer mu.Unlock()

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

	doctorList = append(doctorList, doctorEntry{
		InterfaceName: t.String(),
		DriverName:    provider.Name(),
		Check:         provider.CheckCompatibility,
	})
}

func GetByName[T any](name string) (DriverProvider[T], error) {
	mu.RLock()
	defer mu.RUnlock()

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
	mu.RLock()
	providersCopy := make([]DriverProvider[T], 0)
	var t reflect.Type = reflect.TypeFor[T]()
	if t.Kind() != reflect.Interface {
		mu.RUnlock()
		var zero T
		return zero, ErrNotInterface
	}

	if providers, ok := Drivers[t]; ok {
		for _, p := range providers {
			providersCopy = append(providersCopy, p.(DriverProvider[T]))
		}
	}
	mu.RUnlock()

	var zero T
	if len(providersCopy) == 0 {
		return zero, ErrNotFound
	}

	var report []string

	for _, provider := range providersCopy {
		// 1. Check Compatibility (Doctor)
		if err := provider.CheckCompatibility(ctx); err != nil {
			report = append(report, fmt.Sprintf("❌ [SKIP] %s: %v", provider.Name(), err))
			continue
		}

		// 2. Instantiate
		instance, err := provider.New(ctx)
		if err != nil {
			report = append(report, fmt.Sprintf("⚠️ [FAIL] %s: initialization failed: %v", provider.Name(), err))
			continue
		}

		return instance, nil
	}

	return zero, fmt.Errorf("no available driver for %s:\n%s", t.String(), strings.Join(report, "\n"))
}

type DriverStatus struct {
	Name      string
	Available bool
	Error     error
}

type InterfaceStatus struct {
	Name    string
	Drivers []DriverStatus
}

// Doctor returns the status of all registered drivers
func Doctor(ctx context.Context) []InterfaceStatus {
	mu.RLock()
	defer mu.RUnlock()

	// Group by Type Name
	byType := make(map[string][]doctorEntry)
	for _, d := range doctorList {
		byType[d.InterfaceName] = append(byType[d.InterfaceName], d)
	}

	// Sort types for output stability
	var types []string
	for t := range byType {
		types = append(types, t)
	}
	sort.Strings(types)

	var result []InterfaceStatus

	for _, tName := range types {
		drivers := byType[tName]
		ifaceStatus := InterfaceStatus{
			Name: tName,
		}

		for _, d := range drivers {
			err := d.Check(ctx)
			status := DriverStatus{
				Name:      d.DriverName,
				Available: err == nil,
				Error:     err,
			}
			ifaceStatus.Drivers = append(ifaceStatus.Drivers, status)
		}
		result = append(result, ifaceStatus)
	}

	return result
}
