package tekojar

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	SYSTEM     = "SYSTEM"
	FormatTime = "2006-01-02 15:04:05"
)

var output io.Writer = os.Stdout

func InitLogger() error {
	pathDir, _ := ConcatWithExecutablePath("logs")
	os.MkdirAll(pathDir, 0o755)

	logPath := pathDir + "/tekojar.log"
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	
	output = f
	
	return nil
}

func PrintStringErr(name string, pid int, err string) string {
	return fmt.Sprintf("%s - error - [%s | %d] : %s\n", time.Now().Format(FormatTime), name, pid, err)
}

func PrintErr(name string, pid int, errs ...string) {
	msg := strings.Join(errs, " | ")
	log := fmt.Sprintf("%s - error - [%s | %d] :  %s\n", time.Now().Format(FormatTime), name, pid, msg)
	fmt.Fprintf(output, "%s", log)
	fmt.Printf("%s", log)
}

func PrintLog(name string, pid int, messages ...string) {
	msg := strings.Join(messages, " | ")
	log := fmt.Sprintf("%s - info - [%s | %d] :  %s\n", time.Now().Format(FormatTime), name, pid, msg)
	fmt.Fprintf(output, "%s", log)
	fmt.Printf("%s", log)
}

func ConcatWithExecutablePath(path string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	path = filepath.Join(filepath.Dir(exe), path)
	return path, nil
}
