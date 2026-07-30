package networking

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"sync"
	"unsafe"

	"golang.org/x/sys/unix"
)

type TunDevice interface {
	Create(ctx context.Context, name string, mtu int) error
	Up(ctx context.Context) error
	Down(ctx context.Context) error
	Addresses(ctx context.Context) ([]string, error)
	SetIP(ctx context.Context, ipNet string) error
	SetMTU(ctx context.Context, mtu int) error
	Delete(ctx context.Context) error
}

var (
	ErrTUNNotCreated = fmt.Errorf("tun not created")
	ErrTUNExists     = fmt.Errorf("tun already created")
)

type tunDevice struct {
	mu   sync.RWMutex
	fd   int
	name string
	mtu  int
	up   bool
}

func NewTUN() TunDevice {
	return &tunDevice{fd: -1}
}

func (t *tunDevice) Create(ctx context.Context, name string, mtu int) error {
	if name == "" {
		name = "ghost0"
	}
	if mtu <= 0 {
		mtu = 1420
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd >= 0 {
		return ErrTUNExists
	}

	fd, err := unix.Open("/dev/net/tun", unix.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("open /dev/net/tun: %w", err)
	}

	var ifr [32]byte
	copy(ifr[:], name)
	*(*uint16)(unsafe.Pointer(&ifr[16])) = unix.IFF_TUN | unix.IFF_NO_PI

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TUNSETIFF), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		unix.Close(fd)
		return fmt.Errorf("ioctl TUNSETIFF: %v", errno)
	}

	if err := unix.SetNonblock(fd, true); err != nil {
		unix.Close(fd)
		return fmt.Errorf("set nonblock: %w", err)
	}

	devName := string(ifr[:])
	for i, b := range devName {
		if b == 0 {
			devName = devName[:i]
			break
		}
	}

	t.fd = fd
	t.name = devName
	t.mtu = mtu

	return nil
}

func (t *tunDevice) Up(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd < 0 {
		return ErrTUNNotCreated
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(fd)

	var ifr [32]byte
	copy(ifr[:], t.name)

	*(*uint16)(unsafe.Pointer(&ifr[16])) = unix.IFF_UP | unix.IFF_RUNNING

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFFLAGS), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFFLAGS: %v", errno)
	}

	t.up = true

	var mtuIfr [32]byte
	copy(mtuIfr[:], t.name)
	*(*int32)(unsafe.Pointer(&mtuIfr[16])) = int32(t.mtu)

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFMTU), uintptr(unsafe.Pointer(&mtuIfr)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFMTU: %v", errno)
	}

	return nil
}

func (t *tunDevice) Down(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd < 0 {
		return ErrTUNNotCreated
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(fd)

	var ifr [32]byte
	copy(ifr[:], t.name)

	*(*uint16)(unsafe.Pointer(&ifr[16])) = 0

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFFLAGS), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFFLAGS: %v", errno)
	}

	t.up = false
	return nil
}

func (t *tunDevice) Addresses(ctx context.Context) ([]string, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	if t.fd < 0 || !t.up {
		return nil, ErrTUNNotCreated
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return nil, fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(fd)

	var ifr [32]byte
	copy(ifr[:], t.name)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCGIFADDR), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return nil, nil
	}

	ip := net.IP(ifr[20:24]).String()
	return []string{ip}, nil
}

func (t *tunDevice) SetIP(ctx context.Context, ipNet string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd < 0 {
		return ErrTUNNotCreated
	}

	ip, ipnet, err := net.ParseCIDR(ipNet)
	if err != nil {
		return fmt.Errorf("parse cidr: %w", err)
	}

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(fd)

	var ifr [32]byte
	copy(ifr[:], t.name)

	ip4 := ip.To4()
	copy(ifr[4:8], ip4)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFADDR), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFADDR: %v", errno)
	}

	var mask [32]byte
	copy(mask[:], t.name)
	ones, _ := ipnet.Mask.Size()
	var maskBytes [4]byte
	for i := 0; i < 4; i++ {
		if ones >= 8 {
			maskBytes[i] = 0xff
			ones -= 8
		} else if ones > 0 {
			maskBytes[i] = byte(0xff << (8 - ones))
			ones = 0
		}
	}
	copy(mask[4:8], maskBytes[:])

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFNETMASK), uintptr(unsafe.Pointer(&mask)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFNETMASK: %v", errno)
	}

	return nil
}

func (t *tunDevice) SetMTU(ctx context.Context, mtu int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd < 0 {
		return ErrTUNNotCreated
	}

	t.mtu = mtu

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_DGRAM, 0)
	if err != nil {
		return fmt.Errorf("socket: %w", err)
	}
	defer unix.Close(fd)

	var ifr [32]byte
	copy(ifr[:], t.name)
	*(*int32)(unsafe.Pointer(&ifr[16])) = int32(mtu)

	_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.SIOCSIFMTU), uintptr(unsafe.Pointer(&ifr)))
	if errno != 0 {
		return fmt.Errorf("ioctl SIOCSIFMTU: %v", errno)
	}

	return nil
}

func (t *tunDevice) Delete(ctx context.Context) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.fd < 0 {
		return nil
	}

	unix.Close(t.fd)
	t.fd = -1

	if t.name != "" {
		exec.CommandContext(ctx, "ip", "link", "delete", t.name).Run()
	}

	t.name = ""
	t.up = false
	return nil
}

func (t *tunDevice) Read(buf []byte) (int, error) {
	t.mu.RLock()
	fd := t.fd
	t.mu.RUnlock()

	if fd < 0 {
		return 0, ErrTUNNotCreated
	}

	return unix.Read(fd, buf)
}

func (t *tunDevice) Write(buf []byte) (int, error) {
	t.mu.RLock()
	fd := t.fd
	t.mu.RUnlock()

	if fd < 0 {
		return 0, ErrTUNNotCreated
	}

	return unix.Write(fd, buf)
}
