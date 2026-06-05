package app

import (
	"context"
	"strings"

	"tekotools/backend/tekojar"

	"github.com/jinzhu/copier"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TekojarApp struct {
	ctx     context.Context
	tekojar *tekojar.Tekojar
}

func NewTekojarApp() (*TekojarApp, error) {
	s, err := tekojar.LoadSetting()
	if err != nil {
		return nil, err
	}

	tj, err := tekojar.NewWithSetting(s)
	if err != nil {
		return nil, err
	}
	return &TekojarApp{tekojar: tj}, nil
}

func (ta *TekojarApp) Startup(ctx context.Context) {
	ta.ctx = ctx
}

func (ta *TekojarApp) Shutdown(ctx context.Context) {
	ta.tekojar.StopAll()
}

func (ta *TekojarApp) GetSetting() (DTOTekojarSetting, error) {
	ts := ta.tekojar.Setting.Current()
	dto := DTOTekojarSetting{}
	if err := copier.Copy(&dto, ts); err != nil {
		return dto, err
	}
	return dto, nil
}

func (ta *TekojarApp) SaveSetting(dto DTOTekojarSetting) error {
	ts := tekojar.TekojarSetting{}
	if err := copier.Copy(&ts, &dto); err != nil {
		return err
	}
	if err := ta.tekojar.ApplySetting(ts); err != nil {
		return err
	}
	return nil
}

func (ta *TekojarApp) GetAll() []DTOService {
	services := ta.tekojar.GetAll()

	servicesView := make([]DTOService, 0, len(services))
	for _, s := range services {
		servicesView = append(servicesView, DTOService{
			ID:     s.ID,
			Name:   s.Name,
			Status: string(s.Status),
			Idx:    s.Idx,
			Delay:  s.Delay,
		})
	}

	return servicesView
}

func (ta *TekojarApp) Get(id string) (DTOService, error) {
	s, err := ta.tekojar.GetService(id)
	if err != nil {
		return DTOService{}, err
	}
	return DTOService{
		ID:     s.ID,
		Name:   s.Name,
		Status: string(s.Status),
		Idx:    s.Idx,
		Delay:  s.Delay,
	}, nil
}

func (ta *TekojarApp) Start(id string) error {
	ch, err := ta.tekojar.WatchService(id)
	if err != nil {
		return err
	}

	go func() {
		for l := range ch {
			eventName := "service:log"
			runtime.EventsEmit(ta.ctx, eventName, map[string]interface{}{
				"id": id,
				"logView": DTOLog{
					LogType: string(l.LogType),
					Log:     l.Message,
				},
			})
		}
	}()

	if err := ta.tekojar.Start(id); err != nil {
		return err
	}

	return nil
}

func (ta *TekojarApp) Stop(id string) error {
	if err := ta.tekojar.Stop(id); err != nil {
		return err
	}
	return nil
}

// InitDTO exists to force Wails to generate TypeScript types for all DTO
// it is not meant to be called by the frontend
func (ta *TekojarApp) InitDTO() DTORegistry {
	return DTORegistry{}
}

func (ta *TekojarApp) containsIgnoreCase(str string, char string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(char))
}
