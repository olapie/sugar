package mob

import (
	"code.olapie.com/sugar/v2/xjson"
	"encoding/json"
)

type DeviceInfo struct {
	Name       string `json:"name,omitempty"`
	Model      string `json:"model,omitempty"`
	ModelType  string `json:"model_type,omitempty"`
	Language   string `json:"language,omitempty"`
	SysName    string `json:"sys_name,omitempty"`
	SysVersion string `json:"sys_version,omitempty"`
	Carrier    string `json:"carrier,omitempty"`
}

func NewDeviceInfo() *DeviceInfo {
	return new(DeviceInfo)
}

func (d *DeviceInfo) Attributes() map[string]string {
	m := make(map[string]string)
	err := json.Unmarshal(xjson.ToBytes(d), &m)
	if err != nil {
		panic(err)
	}
	return m
}

type AppInfo struct {
	AppID   string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func NewAppInfo() *AppInfo {
	return new(AppInfo)
}

func (i *AppInfo) Attributes() map[string]string {
	m := make(map[string]string)
	err := json.Unmarshal(xjson.ToBytes(i), &m)
	if err != nil {
		panic(err)
	}
	return m
}
