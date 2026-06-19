package tekojar

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"

	_ "net/http/pprof"

	"github.com/google/uuid"
)

type Tekojar struct {
	Setting *Setting

	services map[string]*Service

	mu sync.RWMutex
}

func New() (*Tekojar, error) {
	initialSetup()

	s, err := LoadSetting()
	if err != nil {
		PrintErr(SYSTEM, 0, err.Error())
		return nil, err
	}
	t := &Tekojar{
		services: make(map[string]*Service),
		Setting:  s,
	}
	t.registerServices()
	return t, nil
}

func NewWithSetting(Setting *Setting) (*Tekojar, error) {
	initialSetup()

	t := &Tekojar{
		services: make(map[string]*Service),
		Setting:  Setting,
	}
	t.registerServices()
	return t, nil
}

func initialSetup() error {
	if err := InitLogger(); err != nil {
		panic(err)
	}
	listenAndServeDebugging()
	return nil
}

func listenAndServeDebugging() {
	PrintLog(SYSTEM, 0, "debug on http://localhost:6060/debug/pprof")
	go http.ListenAndServe(":6060", nil)
}

func (t *Tekojar) registerServices() {
	t.mu.Lock()
	defer t.mu.Unlock()

	ts := t.Setting.Current()
	if ts == nil {
		PrintErr(SYSTEM, 0, "setting not loaded")
	}

	if ts != nil && len(ts.ServiceSettings) == 0 {
		PrintLog(SYSTEM, 0, "no services configured")
	}

	if t.services == nil {
		t.services = make(map[string]*Service)
	}

	t.syncServices(ts)
}

// ApplySetting saves and applies updated settings to the running services
func (t *Tekojar) ApplySetting(ts TekojarSetting) error {
	for i := range ts.ServiceSettings {
		if ts.ServiceSettings[i].ID == "" {
			ts.ServiceSettings[i].ID = uuid.NewString()
		}
	}

	if err := t.Setting.Save(ts); err != nil {
		return fmt.Errorf("failed to save setting: %w", err)
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	t.syncServices(&ts)

	return nil
}

func (t *Tekojar) syncServices(ts *TekojarSetting) {
	existOnSetting := make(map[string]struct{})

	for _, ss := range ts.ServiceSettings {
		if ss.ID == "" {
			continue
		}

		if s, exists := t.services[ss.ID]; exists {
			s.UpdateProcess(ss)
		} else {
			t.services[ss.ID] = InitService(ss)
		}

		existOnSetting[ss.ID] = struct{}{}
	}

	// Remove services no longer in settings
	for id := range t.services {
		if _, exists := existOnSetting[id]; !exists {
			delete(t.services, id)
		}
	}
}

func (t *Tekojar) GetService(id string) (*Service, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	s := t.services[id]
	if s == nil {
		PrintErr(SYSTEM, 0, "Service not found")
		return nil, errors.New("Service not found")
	}

	return t.services[id], nil
}

func (t *Tekojar) WatchService(id string) (chan Log, error) {
	s, err := t.GetService(id)
	if err != nil {
		return nil, err
	}
	return s.Subscribe(), nil
}

func (t *Tekojar) UnwatchService(id string) error {
	s, err := t.GetService(id)
	if err != nil {
		return err
	}
	s.Unsubscribe()
	return nil
}

func (t *Tekojar) Start(id string) error {
	s, err := t.GetService(id)
	if err != nil {
		return err
	}

	ts := t.Setting.Current()

	go s.StartProcess(ts.Command)

	return nil
}

func (t *Tekojar) Stop(id string) error {
	s, err := t.GetService(id)
	if err != nil {
		return err
	}
	s.StopProcess()
	return nil
}

func (t *Tekojar) GetAll() []*Service {
	t.mu.RLock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		services = append(services, v)
	}
	t.mu.RUnlock()

	sort.Slice(services, func(i, j int) bool {
		return services[i].Idx < services[j].Idx
	})

	return services
}

func (t *Tekojar) StartAll() {
	PrintLog(SYSTEM, 0, "starting application...")

	t.mu.RLock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		if v.GetStatus() == INACTIVE {
			services = append(services, v)
		}
	}
	t.mu.RUnlock()

	for _, s := range services {
		t.Start(s.ID)
	}
}

func (t *Tekojar) StopAll() {
	t.mu.RLock()
	services := make([]*Service, 0, len(t.services))
	for _, v := range t.services {
		services = append(services, v)
	}
	t.mu.RUnlock()

	for _, v := range services {
		t.Stop(v.ID)
	}

	actives, inactives := t.GetTotalStatusServices()
	PrintLog(SYSTEM, 0, "application stopped", fmt.Sprintf("actives: %d", actives), fmt.Sprintf("inactives: %d", inactives))
}

func (t *Tekojar) StartAndListen() {
	t.StartAll()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	t.StopAll()
}

func (t *Tekojar) GetTotalStatusServices() (int, int) {
	active := 0
	inactive := 0

	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, s := range t.services {
		switch s.GetStatus() {
		case ACTIVE:
			active++
		case INACTIVE:
			inactive++
		}
	}

	return active, inactive
}
