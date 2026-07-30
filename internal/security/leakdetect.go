package security

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ghoststack/ghoststack/internal/platform/linux"
)

type LeakCheckResult struct {
	DNS      bool   `json:"dns"`
	IPv6     bool   `json:"ipv6"`
	ICMP     bool   `json:"icmp"`
	PublicIP string `json:"public_ip"`
	Error    string `json:"error,omitempty"`
}

func CheckLeaks(ctx context.Context, iface string) (*LeakCheckResult, error) {
	result := &LeakCheckResult{}

	dnsResult := checkDNSLeak(ctx)
	result.DNS = dnsResult

	ipv6Result := checkIPv6Leak(ctx)
	result.IPv6 = ipv6Result

	icmpResult := checkICMPLeak(ctx)
	result.ICMP = icmpResult

	publicIP, err := getPublicIP(ctx)
	if err == nil {
		result.PublicIP = publicIP
	}

	return result, nil
}

func checkDNSLeak(ctx context.Context) bool {
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 3 * time.Second}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	_, err := resolver.LookupHost(ctx, "leaktest.ghoststack.local")
	return err == nil
}

func checkIPv6Leak(ctx context.Context) bool {
	if !linux.Supported() {
		return false
	}

	iface, err := net.InterfaceByName("ghost0")
	if err != nil {
		return false
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return false
	}

	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if ok && ipnet.IP.To4() == nil {
			return true
		}
	}
	return false
}

func checkICMPLeak(ctx context.Context) bool {
	conn, err := net.DialTimeout("ip:icmp", "8.8.8.8", 3*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func getPublicIP(ctx context.Context) (string, error) {
	resolver := net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 5 * time.Second}
			return d.DialContext(ctx, "udp", "1.1.1.1:53")
		},
	}

	ips, err := resolver.LookupHost(ctx, "myip.ghoststack.local")
	if err != nil {
		return "", fmt.Errorf("resolve public ip: %w", err)
	}
	if len(ips) > 0 {
		return ips[0], nil
	}
	return "", fmt.Errorf("no public ip found")
}
