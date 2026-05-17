package tekojar

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Tekojar struct {
	config   *Setting
	services map[string]*Service

	mu sync.RWMutex
}

func New(config *Setting) *Tekojar {
	m := &Tekojar{
		services: make(map[string]*Service),
		config:   config,
	}
	m.registerServices()
	return m
}

func (m *Tekojar) registerServices() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.config.Services == nil {
		panic("no service inputed")
	}

	PrintLog(SYSTEM, 0, fmt.Sprintf("config -> %v", m.config))

	if m.services == nil {
		m.services = make(map[string]*Service)
	}

	for _, s := range m.config.Services {
		if s.SkipFlag {
			continue
		}
		m.services[s.Name] = InitService(s.Name, s.Path)
	}
}

func (m *Tekojar) GetService(name string) *Service {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.services[name]
}

func (m *Tekojar) WatchService(name string) chan string {
	s := m.GetService(name)
	if s == nil {
		return nil
	}
	return s.Subscribe()
}

func (m *Tekojar) UnwatchService(name string) {
	s := m.GetService(name)
	if s == nil {
		return
	}
	s.Unsubscribe()
}

func (m *Tekojar) StartAll() {
	PrintLog(SYSTEM, 0, "starting application...")

	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		if v.GetStatus() == INACTIVE {
			services = append(services, v)
		}
	}
	m.mu.RUnlock()

	for _, s := range services {
		cmd := "cmd"
		if m.config.Command != "" {
			cmd = m.config.Command
		}
		if s.Path != "" && strings.Contains(cmd, "$PATH") {
			cmd = strings.ReplaceAll(cmd, "$PATH", s.Path)
		}

		// split and remove first args, because the first args is command
		args := strings.Split(cmd, " ")
		if len(args) > 1 {
			cmd = args[0]
			args = args[1:]
		}

		go s.StartProcess(cmd, args...)
	}

	m.ListenShutdown()
}

func (m *Tekojar) StopAll() {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		services = append(services, v)
	}
	m.mu.RUnlock()

	for _, v := range services {
		v.StopProcess()
	}

	actives, inactives := m.GetTotalStatusServices()
	PrintLog(SYSTEM, 0, "application stopped", fmt.Sprintf("actives: %d", actives), fmt.Sprintf("inactives: %d", inactives))
}

func (m *Tekojar) GetTotalStatusServices() (int, int) {
	active := 0
	inactive := 0

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, s := range m.services {
		switch s.GetStatus() {
		case ACTIVE:
			active++
		case INACTIVE:
			inactive++
		}
	}

	return active, inactive
}

// ListenShutdown will block goroutines
func (m *Tekojar) ListenShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	if m.config.AutoShutdown {
		m.AutomaticShutDownTicker(5*time.Second, sigChan)
	}

	<-sigChan

	m.StopAll()
}

func (m *Tekojar) AutomaticShutDownTicker(interval time.Duration, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			_, totalInactive := m.GetTotalStatusServices()
			if len(m.services) == totalInactive {
				sigChan <- syscall.SIGQUIT
				return // exit goroutine after triggering
			}
		}
	}()
}
