//go:build windows

package tekojar

import (
	"os/exec"
	"strconv"
	"syscall"
)

func getSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		HideWindow: true,
	}
}

func sendShutdownSignal(s *Service) error {
	// taskkill /T sends a shutdown message to the process tree.
	// Without /F it's not forced — gives Spring Boot a chance to run shutdown hooks.
	cmd := exec.Command("taskkill", "/T", "/PID", strconv.Itoa(s.Pid))
	return cmd.Run()
}
