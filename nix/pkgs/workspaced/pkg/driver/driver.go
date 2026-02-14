package driver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	ID() string   // Unique slug for the driver (e.g. "wayland_swaybg")
	Name() string // Human readable name
	DefaultWeight() int
	CheckCompatibility(ctx context.Context) error
	New(ctx context.Context) (T, error)
}

type doctorEntry struct {
	InterfaceType reflect.Type
	InterfaceName string
	DriverID      string
	DriverName    string
	Check         func(context.Context) error
	DefaultWeight func() int
}

const (
	DefaultWeight = 50
)

var (
	mu            sync.RWMutex
	Drivers       = map[reflect.Type]map[string]any{}
	driverWeights = map[string]map[string]int{}
	doctorList    = []doctorEntry{}
)

// SetWeights configures driver priorities. Weights must be between 0 and 100.
func SetWeights(w map[string]map[string]int) error {
	mu.Lock()
	defer mu.Unlock()

	for iface, drivers := range w {
		for id, weight := range drivers {
			if weight < 0 || weight > 100 {
				return fmt.Errorf("invalid weight %d for driver %q in interface %q: must be between 0 and 100", weight, id, iface)
			}
		}
	}
	driverWeights = w
	return nil
}

func Register[T any](provider DriverProvider[T]) {
	mu.Lock()
	defer mu.Unlock()

	var t reflect.Type = reflect.TypeFor[T]()
	if t.Kind() != reflect.Interface {
		panic(fmt.Errorf("driver %s is not an interface", t.String()))
	}
	id := provider.ID()
	if id == "" {
		panic(fmt.Errorf("driver for %s registered with empty ID", t.String()))
	}

	if _, ok := Drivers[t]; !ok {
		Drivers[t] = make(map[string]any)
	}
	if _, ok := Drivers[t][id]; ok {
		panic(fmt.Errorf("driver ID %q already registered for interface %s", id, t.String()))
	}
	Drivers[t][id] = provider

	doctorList = append(doctorList, doctorEntry{
		InterfaceType: t,
		InterfaceName: getInterfaceName(t),
		DriverID:      id,
		DriverName:    provider.Name(),
		Check:         provider.CheckCompatibility,
		DefaultWeight: provider.DefaultWeight,
	})
}

func getInterfaceName(t reflect.Type) string {
	if t.PkgPath() != "" {
		return t.PkgPath() + "." + t.Name()
	}
	return t.String()
}

func Get[T any](ctx context.Context) (T, error) {
	mu.RLock()
	var t reflect.Type = reflect.TypeFor[T]()
	if t.Kind() != reflect.Interface {
		mu.RUnlock()
		var zero T
		return zero, ErrNotInterface
	}

	ifaceName := getInterfaceName(t)
	weights := driverWeights[ifaceName]

	providers := make([]DriverProvider[T], 0)
	if pMap, ok := Drivers[t]; ok {
		for _, p := range pMap {
			providers = append(providers, p.(DriverProvider[T]))
		}
	}
	mu.RUnlock()

	var zero T
	if len(providers) == 0 {
		return zero, ErrNotFound
	}

	// Sort providers by weight then ID
	sort.Slice(providers, func(i, j int) bool {
		wi := providers[i].DefaultWeight()
		if w, ok := weights[providers[i].ID()]; ok {
			wi = w
		}
		wj := providers[j].DefaultWeight()
		if w, ok := weights[providers[j].ID()]; ok {
			wj = w
		}

		if wi != wj {
			return wi > wj // Higher weight first
		}
		return providers[i].ID() < providers[j].ID() // Deterministic fallback
	})

	var report []string

	for _, provider := range providers {
		weight := provider.DefaultWeight()
		if w, ok := weights[provider.ID()]; ok {
			weight = w
		}

		if err := provider.CheckCompatibility(ctx); err != nil {
			report = append(report, fmt.Sprintf("❌ [SKIP] %s (%s) weight=%d: %v", provider.ID(), provider.Name(), weight, err))
			slog.Debug("driver skipped", "interface", ifaceName, "id", provider.ID(), "name", provider.Name(), "weight", weight, "error", err)
			continue
		}

		instance, err := provider.New(ctx)
		if err != nil {
			report = append(report, fmt.Sprintf("⚠️ [FAIL] %s (%s) weight=%d: initialization failed: %v", provider.ID(), provider.Name(), weight, err))
			slog.Debug("driver init failed", "interface", ifaceName, "id", provider.ID(), "name", provider.Name(), "weight", weight, "error", err)
			continue
		}

		slog.Debug("driver selected", "interface", ifaceName, "id", provider.ID(), "name", provider.Name(), "weight", weight)
		return instance, nil
	}

	return zero, fmt.Errorf("no available driver for %s:\n%s", t.String(), strings.Join(report, "\n"))
}

type DriverStatus struct {
	ID        string
	Name      string
	Weight    int
	Available bool
	Selected  bool
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

	byType := make(map[reflect.Type][]doctorEntry)
	for _, d := range doctorList {
		byType[d.InterfaceType] = append(byType[d.InterfaceType], d)
	}

	var types []reflect.Type
	for t := range byType {
		types = append(types, t)
	}
	sort.Slice(types, func(i, j int) bool {
		return getInterfaceName(types[i]) < getInterfaceName(types[j])
	})

	var result []InterfaceStatus

	for _, t := range types {
		entries := byType[t]
		ifaceName := getInterfaceName(t)
		weights := driverWeights[ifaceName]
		ifaceStatus := InterfaceStatus{
			Name: ifaceName,
		}

		for _, d := range entries {
			err := d.Check(ctx)
			weight := d.DefaultWeight()
			if w, ok := weights[d.DriverID]; ok {
				weight = w
			}
			status := DriverStatus{
				ID:        d.DriverID,
				Name:      d.DriverName,
				Weight:    weight,
				Available: err == nil,
				Error:     err,
			}
			ifaceStatus.Drivers = append(ifaceStatus.Drivers, status)
		}

		// Sort drivers in doctor report
		sort.Slice(ifaceStatus.Drivers, func(i, j int) bool {
			if ifaceStatus.Drivers[i].Weight != ifaceStatus.Drivers[j].Weight {
				return ifaceStatus.Drivers[i].Weight > ifaceStatus.Drivers[j].Weight
			}
			return ifaceStatus.Drivers[i].ID < ifaceStatus.Drivers[j].ID
		})

		// Mark selected candidate
		for i := range ifaceStatus.Drivers {
			if ifaceStatus.Drivers[i].Available {
				ifaceStatus.Drivers[i].Selected = true
				break
			}
		}

		result = append(result, ifaceStatus)
	}

	return result
}
