package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
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

	cmd *exec.Cmd

	mu sync.RWMutex

	LogCh   chan string
	logFile *os.File
	doneCh  chan struct{}
}

func InitService(name string, path string) *Service {
	return &Service{
		Name:   name,
		Path:   path + "/" + name,
		Status: INACTIVE,
		LogCh:  make(chan string, 100),
	}
}

// So StartProcess() is still running in its goroutine the whole time — it's blocked at s.cmd.Wait().
// When SIGTERM causes the JAR to exit, cmd.Wait() unblocks naturally and StartProcess() finishes its own cleanup.
// StopProcess() just sends the signal and waits at <-s.doneCh for StartProcess() to confirm everything is cleaned up before returning.
// They work together — StopProcess() triggers the exit, StartProcess() handles the cleanup.

// StartProcess  →  start, stream logs, wait, cleanup
// StopProcess   →  send signal, wait for StartProcess to confirm done

func (s *Service) StartProcess(command string, args ...string) {
	s.doneCh = make(chan struct{})

	os.MkdirAll("logs", 0o755)
	f, err := os.OpenFile("logs/"+s.Name+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		s.abortError(err)
		return
	}

	s.logFile = f

	s.cmd = exec.Command(command, args...)

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

	s.SetStatus(ACTIVE)
	s.Pid = s.cmd.Process.Pid

	fmt.Printf("Process %s with pid %d started\n", s.Name, s.Pid)

	var wg sync.WaitGroup
	wg.Add(2)
	go s.streamLog(stderr, &wg)
	go s.streamLog(stdout, &wg)

	s.cmd.Wait()
	wg.Wait()

	// StartProcess owns all cleanup
	s.cleanUpService()

	fmt.Printf("Process %s with pid %d exited\n", s.Name, s.Pid)
}

func (s *Service) StopProcess() {
	if s.cmd == nil || s.cmd.Process == nil {
		return
	}

	// send signal first, then wait
	s.cmd.Process.Signal(syscall.SIGTERM)

	select {
	case <-s.doneCh:
		fmt.Printf("Process %s stopped cleanly\n", s.Name)
	case <-time.After(10 * time.Second):
		fmt.Printf("Process %s didn't stop in time, force killing\n", s.Name)
		s.cmd.Process.Kill()
		<-s.doneCh // still wait for cleanup to finish
	}
}

func (s *Service) abortError(err error) {
	s.LogCh <- fmt.Sprintf("[%s] error: %s", s.Name, err.Error())
	s.cleanUpService()
}

func (s *Service) cleanUpService() {
	s.SetStatus(INACTIVE)
	if s.logFile != nil {
		s.logFile.Close()
	}
	close(s.LogCh)
	close(s.doneCh) // signal to StopProcess that cleanup is done
}

func (s *Service) streamLog(logPipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line := fmt.Sprintf("[%s] %s", s.Name, scanner.Text())
		s.LogCh <- line
		s.logFile.WriteString(line + "\n")
	}
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
