package power

import (
	"context"
	"fmt"
	"net"
	"workspaced/pkg/api"
	"workspaced/pkg/config"
	"workspaced/pkg/driver"
	"workspaced/pkg/logging"
)

func Lock(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Lock(ctx)
}

func Reboot(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Reboot(ctx)
}

func Shutdown(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Shutdown(ctx)
}

func Suspend(ctx context.Context) error {
	d, err := driver.Get[Driver](ctx)
	if err != nil {
		return err
	}
	return d.Suspend(ctx)
}

func Wake(ctx context.Context, host string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	hostCfg, ok := cfg.Hosts[host]
	macStr := ""
	if !ok {
		return fmt.Errorf("%w: %s", api.ErrHostNotFound, host)
	} else {
		macStr = hostCfg.MAC
	}

	if macStr == "" {
		return fmt.Errorf("%w: host %s has no MAC address", api.ErrConfigNotFound, host)
	}

	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", api.ErrInvalidAddr, macStr, err)
	}

	packet := make([]byte, 6+16*6)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:(i+1)*6], hwAddr)
	}

	conn, err := net.Dial("udp", "255.255.255.255:9")
	if err != nil {
		return fmt.Errorf("failed to dial UDP broadcast: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			logging.ReportError(ctx, err)
		}
	}()

	_, err = conn.Write(packet)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}

	logging.GetLogger(ctx).Info("sent Wake-on-LAN magic packet", "host", host, "mac", macStr)
	return nil
}
