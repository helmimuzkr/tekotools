//go:build !windows

package tekojar

import (
	"syscall"
)

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func sendShutdownSignal(s *Service) error {
	return s.cmd.Process.Signal(syscall.SIGTERM)
}
