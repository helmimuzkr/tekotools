//go:build !windows

package tekojar

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}

func sendStopSignal(s *Service) {
	s.cmd.Process.Signal(syscall.SIGTERM)
}
