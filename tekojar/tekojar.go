package tekojar

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"
)

type Tekojar struct {
	setting  *Setting
	services map[string]*Service

	mu sync.RWMutex
}

func New(setting *Setting) *Tekojar {
	m := &Tekojar{
		services: make(map[string]*Service),
		setting:  setting,
	}
	m.registerServices()
	return m
}

func (m *Tekojar) registerServices() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.setting.Services == nil {
		panic("no service inputed")
	}

	PrintLog(SYSTEM, 0, fmt.Sprintf("setting -> %v", m.setting))

	if m.services == nil {
		m.services = make(map[string]*Service)
	}

	for _, s := range m.setting.Services {
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
		PrintErr(SYSTEM, 0, "Service not found")
		return nil
	}
	return s.Subscribe()
}

func (m *Tekojar) UnwatchService(name string) {
	s := m.GetService(name)
	if s == nil {
		PrintErr(SYSTEM, 0, "Service not found")
		return
	}
	s.Unsubscribe()
}

func (t *Tekojar) Start(name string) {
	s := t.GetService(name)
	if s == nil {
		PrintErr(SYSTEM, 0, "Service not found")
		return
	}

	go s.StartProcess(t.setting.Command)

	go t.ListenShutdown(name)
}

func (t *Tekojar) Stop(name string) {
	s := t.GetService(name)
	if s == nil {
		PrintErr(SYSTEM, 0, "Service not found")
		return
	}
	s.StopProcess()
}

func (m *Tekojar) GetAll() []*Service {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		services = append(services, v)
	}
	m.mu.RUnlock()

	sort.Slice(services, func(i, j int) bool {
		return services[i].Name < services[j].Name
	})

	return services
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
		go s.StartProcess(m.setting.Command)
	}

	m.ListenAndShutdownAll()
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
func (m *Tekojar) ListenAndShutdownAll() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	if m.setting.AutoShutdown {
		m.AutomaticShutDownTicker("", 5*time.Second, sigChan)
	}

	<-sigChan

	m.StopAll()
}

// ListenShutdown will block goroutines
func (m *Tekojar) ListenShutdown(name string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)

	if m.setting.AutoShutdown {
		m.AutomaticShutDownTicker(name, 5*time.Second, sigChan)
	}

	<-sigChan

	m.Stop(name)
}

func (m *Tekojar) AutomaticShutDownTicker(name string, interval time.Duration, sigChan chan os.Signal) {
	ticker := time.NewTicker(interval)

	go func() {
		defer func() {
			ticker.Stop()
			PrintLog(SYSTEM, 0, "Shutdown automaticlly")
		}()
		for range ticker.C {
			_, totalInactive := m.GetTotalStatusServices()
			if name != "" && m.GetService(name).Status == INACTIVE {
				sigChan <- syscall.SIGTERM
				return // exit goroutine after triggering
			} else if len(m.services) == totalInactive {
				sigChan <- syscall.SIGTERM
				return // exit goroutine after triggering
			}
		}
	}()
}
