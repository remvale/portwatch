package portscanner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PortEntry represents a single bound port entry on the system.
type PortEntry struct {
	Protocol string
	LocalAddress string
	LocalPort int
	PID int
	State string
}

// Scanner reads active port bindings from the system.
type Scanner struct{}

// NewScanner creates a new Scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Scan reads /proc/net/tcp and /proc/net/tcp6 and returns active port entries.
func (s *Scanner) Scan() ([]PortEntry, error) {
	var entries []PortEntry

	for _, path := range []string{"/proc/net/tcp", "/proc/net/tcp6"} {
		results, err := parseNetFile(path, "tcp")
		if err != nil {
			continue // file may not exist on all systems
		}
		entries = append(entries, results...)
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no port data available; ensure running on Linux with /proc access")
	}

	return entries, nil
}

func parseNetFile(path, protocol string) ([]PortEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []PortEntry
	scanner := bufio.NewScanner(f)
	scanner.Scan() // skip header line

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) < 12 {
			continue
		}

		localAddr := fields[1]
		state := fields[3]

		port, err := parseHexPort(localAddr)
		if err != nil {
			continue
		}

		addrPart := parseHexAddr(localAddr)

		entries = append(entries, PortEntry{
			Protocol:     protocol,
			LocalAddress: addrPart,
			LocalPort:    port,
			State:        state,
		})
	}

	return entries, scanner.Err()
}

func parseHexPort(addrPort string) (int, error) {
	parts := strings.Split(addrPort, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid address format: %s", addrPort)
	}
	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return 0, err
	}
	return int(port), nil
}

func parseHexAddr(addrPort string) string {
	parts := strings.Split(addrPort, ":")
	if len(parts) != 2 {
		return "unknown"
	}
	return parts[0]
}
