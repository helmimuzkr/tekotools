//go:build windows

package tekojar

import "syscall"

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func sendStopSignal(s *Service) {
	s.cmd.Process.Kill()
}
