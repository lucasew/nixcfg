package power

import (
	"context"
	"fmt"
	"net"
	"workspaced/pkg/config"
	"workspaced/pkg/logging"
)

func Wake(ctx context.Context, host string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	hostCfg, ok := cfg.Hosts[host]
	macStr := ""
	if !ok {
		return fmt.Errorf("host %s not found in config", host)
	} else {
		macStr = hostCfg.MAC
	}

	if macStr == "" {
		return fmt.Errorf("host %s has no MAC address configured", host)
	}

	hwAddr, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("invalid MAC address for %s: %w", host, err)
	}

	// WoL magic packet: 6 bytes of 0xFF followed by 16 repetitions of the MAC address
	packet := make([]byte, 6+16*6)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 1; i <= 16; i++ {
		copy(packet[i*6:(i+1)*6], hwAddr)
	}

	// Send to broadcast address
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
