package gen1

import (
	"fmt"
	"net/http"
	"time"
)

type BaseShellyResponse struct {
	Type         string `json:"type"`
	Mac          string `json:"mac"`
	Auth         bool   `json:"auth"`
	Fw           string `json:"fw"`
	LongId       int    `json:"longid"`
	Discoverable bool   `json:"discoverable"`
}

type BaseSettingsResponse struct {
	Device                    BaseDevice    `json:"device,omitempty"`
	WifiAp                    BaseWifiAp    `json:"wifi_ap,omitempty"`
	WifiSta                   BaseWifiSta   `json:"wifi_sta,omitempty"`
	WifiSta1                  BaseWifiSta1  `json:"wifi_sta1,omitempty"`
	ApRoaming                 BaseApRoaming `json:"ap_roaming,omitempty"`
	Mqtt                      BaseMqtt      `json:"mqtt,omitempty"`
	Coiot                     BaseCoiot     `json:"coiot,omitempty"`
	Sntp                      BaseSntp      `json:"sntp,omitempty"`
	Login                     BaseLogin     `json:"login,omitempty"`
	PinCode                   string        `json:"pin_code,omitempty"`
	Name                      string        `json:"name,omitempty"`
	Fw                        string        `json:"fw,omitempty"`
	Discoverable              bool          `json:"discoverable,omitempty"`
	BuildInfo                 BaseBuildInfo `json:"build_info,omitempty"`
	Cloud                     BaseCloud     `json:"cloud,omitempty"`
	Timezone                  string        `json:"timezone,omitempty"`
	Lat                       float64       `json:"lat,omitempty"`
	Lng                       float64       `json:"lng,omitempty"`
	Tzautodetect              bool          `json:"tzautodetect,omitempty"`
	TzUtcOffset               int           `json:"tz_utc_offset,omitempty"`
	TzDst                     bool          `json:"tz_dst,omitempty"`
	TzDstAuto                 bool          `json:"tz_dst_auto,omitempty"`
	Time                      string        `json:"time,omitempty"`
	Unixtime                  int           `json:"unixtime,omitempty"`
	LedStatusDisable          bool          `json:"led_status_disable,omitempty"`
	DebugEnable               bool          `json:"debug_enable,omitempty"`
	AllowCrossOrigin          bool          `json:"allow_cross_origin,omitempty"`
	WifirecoveryRebootEnabled bool          `json:"wifirecovery_reboot_enabled,omitempty"`
}
type BaseDevice struct {
	Type     string `json:"type,omitempty"`
	Mac      string `json:"mac,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}
type BaseWifiAp struct {
	Enabled bool   `json:"enabled,omitempty"`
	Ssid    string `json:"ssid,omitempty"`
	Key     string `json:"key,omitempty"`
}
type BaseWifiSta struct {
	Enabled    bool   `json:"enabled,omitempty"`
	Ssid       string `json:"ssid,omitempty"`
	Ipv4Method string `json:"ipv4_method,omitempty"`
	IP         any    `json:"ip,omitempty"`
	Gw         any    `json:"gw,omitempty"`
	Mask       any    `json:"mask,omitempty"`
	DNS        any    `json:"dns,omitempty"`
}
type (
	BaseWifiSta1  struct{}
	BaseApRoaming struct {
		Enabled   bool `json:"enabled,omitempty"`
		Threshold int  `json:"threshold,omitempty"`
	}
)

type BaseMqtt struct {
	Enable              bool    `json:"enable,omitempty"`
	Server              string  `json:"server,omitempty"`
	User                string  `json:"user,omitempty"`
	ID                  string  `json:"id,omitempty"`
	ReconnectTimeoutMax float32 `json:"reconnect_timeout_max,omitempty"`
	ReconnectTimeoutMin float32 `json:"reconnect_timeout_min,omitempty"`
	CleanSession        bool    `json:"clean_session,omitempty"`
	KeepAlive           int     `json:"keep_alive,omitempty"`
	MaxQos              int     `json:"max_qos,omitempty"`
	Retain              bool    `json:"retain,omitempty"`
	UpdatePeriod        int     `json:"update_period,omitempty"`
}
type BaseCoiot struct {
	Enabled      bool   `json:"enabled,omitempty"`
	UpdatePeriod int    `json:"update_period,omitempty"`
	Peer         string `json:"peer,omitempty"`
}
type BaseSntp struct {
	Server  string `json:"server,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}
type BaseLogin struct {
	Enabled     bool   `json:"enabled,omitempty"`
	Unprotected bool   `json:"unprotected,omitempty"`
	Username    string `json:"username,omitempty"`
	Password    string `json:"-"`
}
type BaseBuildInfo struct {
	BuildID        string    `json:"build_id,omitempty"`
	BuildTimestamp time.Time `json:"build_timestamp,omitempty"`
	BuildVersion   string    `json:"build_version,omitempty"`
}
type BaseCloud struct {
	Enabled   bool `json:"enabled,omitempty"`
	Connected bool `json:"connected,omitempty"`
}

type BaseOtaResponse struct {
	Status      string `json:"status"`
	HasUpdate   bool   `json:"has_update"`
	NewVersion  string `json:"new_version"`
	OldVersion  string `json:"old_version"`
	BetaVersion string `json:"beta_version"`
}

type BaseOtaRequest struct {
	Url    string `json:"url"`
	Update bool   `json:"update"`
	Beta   bool   `json:"beta"`
}

type BaseOtaCheck struct {
	Status string `json:"status"`
}

type BaseWifiScan struct {
	Wifiscan string                `json:"wifiscan,omitempty"`
	Results  []BaseWifiScanResults `json:"results,omitempty"`
}
type BaseWifiScanResults struct {
	Ssid    string `json:"ssid,omitempty"`
	Auth    int    `json:"auth,omitempty"`
	Channel int    `json:"channel,omitempty"`
	Bssid   string `json:"bssid,omitempty"`
	Rssi    int    `json:"rssi,omitempty"`
}

type ShellyService struct {
	client *Client
}

func (s *ShellyService) GetShelly() (*BaseShellyResponse, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/shelly", nil)
	if err != nil {
		return nil, err
	}
	var info BaseShellyResponse
	_, err = s.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (s *ShellyService) GetSettings() (*BaseSettingsResponse, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/settings", nil)
	if err != nil {
		return nil, err
	}
	var info BaseSettingsResponse
	resp, err := s.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	return &info, nil
}

func (s *ShellyService) GetOta() (*BaseOtaResponse, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/ota", nil)
	if err != nil {
		return nil, err
	}
	var info BaseOtaResponse
	resp, err := s.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	return &info, nil
}

func (s *ShellyService) GetOtaCheck() (*BaseOtaCheck, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/ota/check", nil)
	if err != nil {
		return nil, err
	}
	var info BaseOtaCheck
	resp, err := s.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	return &info, nil
}

func (s *ShellyService) GetWifiScan() (*BaseWifiScan, error) {
	req, err := s.client.NewRequest(http.MethodGet, "/wifiscan", nil)
	if err != nil {
		return nil, err
	}
	var info BaseWifiScan
	resp, err := s.client.Do(req, &info)
	if err != nil {
		return nil, err
	}
	fmt.Print(resp)
	return &info, nil
}
