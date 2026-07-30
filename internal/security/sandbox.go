package security

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type SandboxPolicy struct {
	AllowNetwork   bool
	AllowFilesystem bool
	AllowedPaths   []string
	AllowExec      bool
	DeniedSyscalls []string
}

type Sandbox interface {
	Apply(ctx context.Context, policy SandboxPolicy) error
	ApplyToProcess(pid int, policy SandboxPolicy) error
}

type sandboxImpl struct{}

func NewSandbox() Sandbox {
	return &sandboxImpl{}
}

func (s *sandboxImpl) Apply(ctx context.Context, policy SandboxPolicy) error {
	return s.applySeccomp(policy)
}

func (s *sandboxImpl) ApplyToProcess(pid int, policy SandboxPolicy) error {
	return s.applyToProcess(pid, policy)
}

func (s *sandboxImpl) applySeccomp(policy SandboxPolicy) error {
	whitelist := []uintptr{
		syscall.SYS_READ,
		syscall.SYS_WRITE,
		syscall.SYS_CLOSE,
		syscall.SYS_MMAP,
		syscall.SYS_MUNMAP,
		syscall.SYS_MPROTECT,
		syscall.SYS_BRK,
		syscall.SYS_EXIT,
		syscall.SYS_EXIT_GROUP,
		syscall.SYS_NANOSLEEP,
		syscall.SYS_GETTIMEOFDAY,
		syscall.SYS_CLOCK_GETTIME,
		syscall.SYS_GETPID,
		syscall.SYS_GETTID,
		syscall.SYS_TGKILL,
		syscall.SYS_FUTEX,
		syscall.SYS_OPENAT,
		syscall.SYS_NEWFSTATAT,
		syscall.SYS_READLINKAT,
		syscall.SYS_FSTAT,
		syscall.SYS_LSEEK,
		syscall.SYS_WRITEV,
		syscall.SYS_PREAD64,
		syscall.SYS_PWRITE64,
		syscall.SYS_GETDENTS64,
		syscall.SYS_SCHED_YIELD,
		syscall.SYS_RT_SIGACTION,
		syscall.SYS_RT_SIGPROCMASK,
		syscall.SYS_RT_SIGRETURN,
		syscall.SYS_SIGALTSTACK,
	}

	if policy.AllowNetwork {
		whitelist = append(whitelist,
			syscall.SYS_SOCKET,
			syscall.SYS_CONNECT,
			syscall.SYS_BIND,
			syscall.SYS_LISTEN,
			syscall.SYS_ACCEPT,
			syscall.SYS_SENDTO,
			syscall.SYS_RECVFROM,
			syscall.SYS_SENDMSG,
			syscall.SYS_RECVMSG,
			syscall.SYS_GETSOCKOPT,
			syscall.SYS_SETSOCKOPT,
			syscall.SYS_GETSOCKNAME,
			syscall.SYS_SHUTDOWN,
		)
	}

	if policy.AllowExec {
		whitelist = append(whitelist,
			syscall.SYS_EXECVE,
			syscall.SYS_EXECVEAT,
			syscall.SYS_CLONE,
			syscall.SYS_FORK,
			syscall.SYS_VFORK,
			syscall.SYS_WAIT4,
			syscall.SYS_WAITID,
		)
	}

	if policy.AllowFilesystem {
		whitelist = append(whitelist,
			syscall.SYS_MKDIRAT,
			syscall.SYS_UNLINKAT,
			syscall.SYS_SYMLINKAT,
			syscall.SYS_RENAMEAT,
			syscall.SYS_LINKAT,
			syscall.SYS_FCHMODAT,
			syscall.SYS_FCHOWNAT,
			syscall.SYS_UTIMENSAT,
		)
	}

	filter := buildSeccompFilter(whitelist)

	prog := &unix.SockFprog{
		Filter: filter,
	}

	if err := unix.Prctl(unix.PR_SET_NO_NEW_PRIVS, 1, 0, 0, 0); err != nil {
		return fmt.Errorf("PR_SET_NO_NEW_PRIVS: %w", err)
	}

	if err := unix.Prctl(unix.PR_SET_SECCOMP, unix.SECCOMP_MODE_FILTER, uintptr(unsafePtr(prog)), 0, 0); err != nil {
		return fmt.Errorf("PR_SET_SECCOMP: %w", err)
	}

	return nil
}

type sockFilter struct {
	Code uint16
	Jt   uint8
	Jf   uint8
	K    uint32
}

func buildSeccompFilter(whitelist []uintptr) []unix.SockFilter {
	var filters []unix.SockFilter

	filters = append(filters, unix.SockFilter{
		Code: 0x20,
		Jt:   0,
		Jf:   0,
		K:    4,
	})

	allowed := make(map[uintptr]bool)
	for _, sysno := range whitelist {
		allowed[sysno] = true
	}

	filters = append(filters, unix.SockFilter{
		Code: 0x15,
		Jt:   0,
		Jf:   0,
		K:    0x40000003,
	})

	filters = append(filters, unix.SockFilter{
		Code: 0x15,
		Jt:   0,
		Jf:   1,
		K:    0,
	})

	filters = append(filters, unix.SockFilter{
		Code: 0x06,
		Jt:   0,
		Jf:   0,
		K:    0x7FFF0000,
	})

	for sysno := range allowed {
		filters = append(filters, unix.SockFilter{
			Code: 0x15,
			Jt:   0,
			Jf:   0,
			K:    uint32(sysno),
		})

		filters = append(filters, unix.SockFilter{
			Code: 0x06,
			Jt:   0,
			Jf:   0,
			K:    0x7FFF0000,
		})
	}

	filters = append(filters, unix.SockFilter{
		Code: 0x06,
		Jt:   0,
		Jf:   0,
		K:    0x00050000,
	})

	return filters
}

func unsafePtr(p *unix.SockFprog) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func (s *sandboxImpl) applyToProcess(pid int, policy SandboxPolicy) error {
	if !policy.AllowNetwork {
		if err := applyNetworkIsolation(pid); err != nil {
			return fmt.Errorf("network isolation: %w", err)
		}
	}

	return nil
}

func applyNetworkIsolation(pid int) error {
	nsPath := fmt.Sprintf("/proc/%d/ns/net", pid)
	netnsFd, err := unix.Open(nsPath, unix.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open net ns: %w", err)
	}
	defer unix.Close(netnsFd)

	newNsFd, err := unix.Open("/proc/1/ns/net", unix.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("open default net ns: %w", err)
	}
	defer unix.Close(newNsFd)

	if err := unix.Setns(netnsFd, unix.CLONE_NEWNET); err != nil {
		return fmt.Errorf("join target net ns: %w", err)
	}

	return nil
}
