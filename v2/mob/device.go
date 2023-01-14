package mob

type DeviceInfo struct {
	Name       string `json:"name"`
	Model      string `json:"model"`
	ModelType  string `json:"model_type"`
	Language   string `json:"language"`
	SysVersion string `json:"sys_version"`
	Carrier    string `json:"carrier"`

	AppVersion string `json:"app_version"`
}

func NewDeviceInfo() *DeviceInfo {
	return new(DeviceInfo)
}
