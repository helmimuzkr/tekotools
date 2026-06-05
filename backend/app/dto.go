package app

type DTORegistry struct {
	Service        DTOService
	ServiceLog     DTOServiceLog
	Log            DTOLog
	ServiceSetting DTOServiceSetting
	TekojarSetting DTOTekojarSetting
}

type DTOService struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Idx    int    `json:"idx"`
	Delay  int    `json:"delay"`
}

type DTOServiceLog struct {
	ServiceID string   `json:"service_id"`
	Logs      []DTOLog `json:"logs"`
}

type DTOLog struct {
	LogType string `json:"log_type"`
	Log     string `json:"log"`
}

type DTOTekojarSetting struct {
	Command         string              `json:"command"`
	AutoShutdown    bool                `json:"auto_shutdown"`
	ServiceSettings []DTOServiceSetting `json:"service_settings"`
}

type DTOServiceSetting struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Path     string `json:"path"`
	SkipFlag bool   `json:"skip_flag"`
	Delay    int    `json:"delay"`
	Idx      int    `json:"idx"`
}
