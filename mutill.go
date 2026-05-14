package main

import (
	"sync"
)

type Mutill struct {
	config   *MutillConfig
	services map[string]*Service

	mu sync.RWMutex

	logStarted bool
	logCh      chan string
}

func (m *Mutill) RegisterService(config *MutillConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if config == nil {
		return
	}

	m.config = config

	if m.config.Services == nil {
		panic("no service inputed")
	}

	if m.services == nil {
		m.services = make(map[string]*Service)
	}

	for _, s := range m.config.Services {
		if s.IsSkip {
			continue
		}
		m.services[s.Name] = InitService(s.Name, s.Path)
	}
}

func (m *Mutill) GetService(name string) *Service {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.services[name]
}

func (m *Mutill) StartAll() {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		if v.Status == INACTIVE {
			services = append(services, v)
		}
	}
	m.mu.RUnlock()
	m.listenLog()

	for _, v := range services {
		cmd := "cmd"
		if m.config.Command != "" {
			cmd = m.config.Command
		}

		args := []string{"/c", v.Path}
		if m.config.Args != nil {
			args = m.config.Args
		}

		go v.StartProcess(cmd, args...)
	}
}

func (m *Mutill) StopAll() {
	m.mu.RLock()
	services := make([]*Service, 0, len(m.services))
	for _, v := range m.services {
		if v.Status == ACTIVE {
			services = append(services, v)
		}
	}
	m.mu.RUnlock()

	for _, v := range services {
		v.StopProcess()
	}
}

func (m *Mutill) listenLog() {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.logStarted {
		panic("ServiceChannelLog() already called")
	}

	m.logStarted = true

	m.logCh = make(chan string, 100*len(m.services))

	var wg sync.WaitGroup

	for _, s := range m.services {
		wg.Add(1)
		go func(svc *Service) {
			defer wg.Done()
			for log := range svc.LogCh {
				m.logCh <- log
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(m.logCh)
	}()
}

func (m *Mutill) ReadLog() chan string {
	if m.logCh == nil {
		m.logCh = make(chan string, 100*len(m.services))
	}
	return m.logCh
}
