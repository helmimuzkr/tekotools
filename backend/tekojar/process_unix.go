//go:build !windows

package tekojar

import "syscall"

func hideWindowAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{}
}
