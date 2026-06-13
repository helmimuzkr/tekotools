package tekojar

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type (
	Status  string
	LogType string
)

const (
	ACTIVE   Status = "ACTIVE"
	INACTIVE Status = "INACTIVE"
	STOPPING Status = "STOPPING"

	INFO  LogType = "INFO"
	ERROR LogType = "ERROR"
	TIMER LogType = "TIMER"
)

type Log struct {
	LogType LogType
	Message string
}

type Service struct {
	ID     string
	Name   string
	Path   string
	Status Status
	IsSkip bool
	Idx    int
	Delay  int

	Pid     int
	cmd     *exec.Cmd
	cmdName string
	args    []string

	mu sync.RWMutex

	processDoneCh chan struct{}

	eventMu    sync.RWMutex
	eventLogCh chan Log
}

func InitService(ss ServiceSetting) *Service {
	return &Service{
		ID:     ss.ID,
		Name:   ss.Name,
		Path:   ss.Path + "/" + ss.Name,
		Status: INACTIVE,
		IsSkip: ss.SkipFlag,
		Idx:    ss.Idx,
		Delay:  ss.Delay,
	}
}

func (s *Service) UpdateProcess(ss ServiceSetting) {
	s.Name = ss.Name
	s.Path = ss.Path + "/" + ss.Name
	s.IsSkip = ss.SkipFlag
	s.Idx = ss.Idx
	s.Delay = ss.Delay
}

// So StartProcess() is still running in its goroutine the whole time — it's blocked at s.cmd.Wait().
// When SIGTERM causes the JAR to exit, cmd.Wait() unblocks naturally and StartProcess() finishes its own cleanup.
// StopProcess() just sends the signal and waits at <-s.processDoneCh for StartProcess() to confirm everything is cleaned up before returning.
// They work together — StopProcess() triggers the exit, StartProcess() handles the cleanup.

// StartProcess  →  start, stream logs, wait, cleanup
// StopProcess   →  send signal, wait for StartProcess to confirm done

func (s *Service) StartProcess(command string) error {
	if s.Status == "ACTIVE" {
		PrintErr(s.Name, s.Pid, "process already active")
		return nil
	}

	s.processDoneCh = make(chan struct{})
	s.SetStatus(ACTIVE)

	s.extractCommand(command)
	s.cmd = exec.Command(s.cmdName, s.args...)
	s.cmd.SysProcAttr = getSysProcAttr()

	PrintLog(s.Name, 0, fmt.Sprintf("executing command : %s %v", s.cmdName, s.args))

	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		s.abortError(err)
		return err
	}
	defer stdout.Close()

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		s.abortError(err)
		return err
	}
	defer stderr.Close()

	if !s.delayTicker() {
		s.waitAndCleanUp(nil)
		return nil
	}

	if err := s.cmd.Start(); err != nil {
		s.abortError(err)
		return err
	}

	s.Pid = s.cmd.Process.Pid
	PrintLog(s.Name, s.Pid, "service started.")

	var wg sync.WaitGroup
	wg.Add(2)
	go s.streamLog(stderr, &wg)
	go s.streamLog(stdout, &wg)

	s.cmd.Wait()
	s.waitAndCleanUp(&wg)
	PrintLog(s.Name, s.Pid, "service exited.")

	return nil
}

func (s *Service) StopProcess() {
	if s.GetStatus() == INACTIVE {
		PrintLog(s.Name, s.Pid, "cannot stop process because process already inactive")
		return
	}

	s.SetStatus(STOPPING)

	if s.cmd != nil && s.cmd.Process != nil {
		// send signal first, then wait
		PrintLog(s.Name, s.Pid, "send shutdown signal to process")
		sendShutdownSignal(s)
	}

	select {
	case <-s.processDoneCh:
		PrintLog(s.Name, s.Pid, "stopped cleanly")
	case <-time.After(10 * time.Second):
		PrintLog(s.Name, s.Pid, "didn't stop in time, force killing")
		s.cmd.Process.Kill()
		<-s.processDoneCh // still wait for cleanup to finish
		PrintLog(s.Name, s.Pid, "successfully force close")
	}
}

func (s *Service) Subscribe() chan Log {
	s.eventMu.Lock()
	defer s.eventMu.Unlock()
	s.eventLogCh = make(chan Log)
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

func (s *Service) abortError(err error) {
	PrintErr(s.Name, s.Pid, err.Error())
	s.waitAndCleanUp(nil)
}

func (s *Service) waitAndCleanUp(wg *sync.WaitGroup) {
	if wg != nil {
		wg.Wait()
	}

	PrintLog(s.Name, s.Pid, "cleaning up ...")
	defer PrintLog(s.Name, s.Pid, "done cleaning up")

	s.eventMu.Lock()
	if s.eventLogCh != nil {
		close(s.eventLogCh)
		s.eventLogCh = nil
	}
	s.eventMu.Unlock()

	if s.processDoneCh != nil {
		close(s.processDoneCh)
	}

	s.SetStatus(INACTIVE)
}

func (s *Service) streamLog(logPipe io.ReadCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	scanner := bufio.NewScanner(logPipe)
	for scanner.Scan() {
		line := scanner.Text()
		s.publishEvent("", line)
	}
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
		s.cmdName = args[0]
		s.args = args[1:]
	} else {
		s.cmdName = cmdName
		s.args = []string{}
	}
}

// copy the channel reference into a local variable, release the lock immediately.
// it because dont want to hold the lock while sending.
// that would cause a deadlock if Unsubscribe tries to acquire the write lock while blocked on a full channel.
func (s *Service) publishEvent(logType LogType, line string) {
	s.eventMu.RLock()
	ch := s.eventLogCh
	s.eventMu.RUnlock()

	log := Log{LogType: logType, Message: line}

	if logType == "" {
		if ContainsIgnoreCase(line, string(ERROR)) {
			log.LogType = ERROR
		} else {
			log.LogType = INFO
		}
	}

	if ch != nil {
		select {
		case ch <- log:
		default: // drop if full
		}
	}
}

func (s *Service) delayTicker() bool {
	ticker := time.NewTicker(time.Second)
	defer func() {
		PrintLog(s.Name, s.Pid, "ticker stopped")
		ticker.Stop()
	}()

	for remaining := s.Delay; remaining > 0; remaining-- {
		if s.GetStatus() != ACTIVE {
			return false
		}
		msg := fmt.Sprintf("starting in %d seconds", remaining)
		PrintLog(s.Name, s.Pid, msg)
		s.publishEvent(TIMER, msg)
		<-ticker.C
	}
	return true
}
