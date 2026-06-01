package app

type DTOService struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Status string   `json:"status"`
	Idx    int      `json:"idx"`
	Delay  int      `json:"delay"`
	Logs   []DTOLog `json:"logs"`
}

type DTOLog struct {
	IsError bool   `json:"is_error"`
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
