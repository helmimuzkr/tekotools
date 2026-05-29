package tekojar

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Status string

const (
	ACTIVE   Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
)

type Service struct {
	Name   string
	Path   string
	Status Status
	Pid    int

	cmd     *exec.Cmd
	cmdName string
	args    []string

	mu sync.RWMutex

	processDoneCh chan struct{}

	logFile *os.File

	eventMu    sync.RWMutex
	eventLogCh chan string
}

func InitService(name string, path string) *Service {
	return &Service{
		Name:   name,
		Path:   path + "/" + name,
		Status: INACTIVE,
	}
}

// So StartProcess() is still running in its goroutine the whole time — it's blocked at s.cmd.Wait().
// When SIGTERM causes the JAR to exit, cmd.Wait() unblocks naturally and StartProcess() finishes its own cleanup.
// StopProcess() just sends the signal and waits at <-s.processDoneCh for StartProcess() to confirm everything is cleaned up before returning.
// They work together — StopProcess() triggers the exit, StartProcess() handles the cleanup.

// StartProcess  →  start, stream logs, wait, cleanup
// StopProcess   →  send signal, wait for StartProcess to confirm done

func (s *Service) StartProcess(command string) {
	if s.Status == "ACTIVE" {
		PrintErr(s.Name, s.Pid, "process already active")
		return
	}

	s.processDoneCh = make(chan struct{})
	s.SetStatus(ACTIVE)

	os.MkdirAll("logs", 0o755)
	f, err := os.OpenFile("logs/"+s.Name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		s.abortError(err)
		return
	}

	s.logFile = f

	s.extractCommand(command)
	s.cmd = exec.Command(s.cmdName, s.args...)

	PrintLog(SYSTEM, 0, fmt.Sprintf("executing command for %s: %s %v", s.Name, s.cmdName, s.args))

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		s.abortError(err)
		return
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		s.abortError(err)
		return
	}

	if err := s.cmd.Start(); err != nil {
		s.abortError(err)
		return
	}

	PrintLog(SYSTEM, 0, fmt.Sprintf("%s started", s.Name))
	s.Pid = s.cmd.Process.Pid

	var wg sync.WaitGroup
	wg.Add(2)
	go s.streamLog(stderr, &wg)
	go s.streamLog(stdout, &wg)

	s.cmd.Wait()
	wg.Wait()

	// StartProcess owns all cleanup
	s.processDone()
	PrintLog(SYSTEM, 0, fmt.Sprintf("%s exited", s.Name))
}

func (s *Service) StopProcess() {
	if s.cmd == nil || s.cmd.Process == nil {
		return
	}

	// send signal first, then wait
	s.cmd.Process.Signal(syscall.SIGTERM)

	select {
	case <-s.processDoneCh:
		PrintLog(s.Name, s.Pid, "stopped cleanly")
	case <-time.After(10 * time.Second):
		PrintLog(s.Name, s.Pid, "sdidn't stop in time, force killing")
		s.cmd.Process.Kill()
		<-s.processDoneCh // still wait for cleanup to finish
	}
}

func (s *Service) abortError(err error) {
	PrintErr(s.Name, s.Pid, err.Error())
	s.processDone()
}

func (s *Service) processDone() {
	PrintLog(SYSTEM, 0, fmt.Sprintf("cleaning up %s ...", s.Name))
	s.SetStatus(INACTIVE)

	if s.logFile != nil {
		s.logFile.Close()
	}

	s.eventMu.Lock()
	if s.eventLogCh != nil {
		close(s.eventLogCh)
		s.eventLogCh = nil
	}
	s.eventMu.Unlock()

	close(s.processDoneCh)
}

func (s *Service) streamLog(logPipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line := scanner.Text()

		// publish log for Subscriber
		s.eventMu.RLock()
		if s.eventLogCh != nil {
			s.eventLogCh <- line // UI subscriber if active
		}
		s.eventMu.RUnlock()

		// file log
		line = fmt.Sprintf("[%s | %d] : %s", s.Name, s.Pid, scanner.Text())
		s.logFile.WriteString(line + "\n")
	}
}

func (s *Service) Subscribe() chan string {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.eventLogCh = make(chan string, 100)
	return s.eventLogCh
}

func (s *Service) Unsubscribe() {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.eventLogCh = nil // stop routing, don't close
}

func (s *Service) SetStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = status
}

func (s *Service) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status
}

func (s *Service) extractCommand(command string) {
	cmdName := "cmd"
	if command != "" {
		cmdName = command
	}
	if s.Path != "" && strings.Contains(cmdName, "$PATH") {
		cmdName = strings.ReplaceAll(cmdName, "$PATH", s.Path)
	}

	// split and remove first args, because the first args is command
	args := strings.Split(cmdName, " ")
	if len(args) > 1 {
		cmdName = args[0]
		args = args[1:]
	}

	s.cmdName = cmdName
	s.args = args
}
