package hetznerrobot

// https://robot.your-server.de/doc/webservice/en.html#boot-configuration

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

type BootProfile struct {
	ActiveProfile   string // linux/rescue/...
	Architecture    string
	AuthorizedKeys  []string
	HostKeys        []string
	Language        string
	OperatingSystem string
	Password        string
	ServerID        int
	ServerIPv4      string
	ServerIPv6      string
}

func (c *HetznerRobotClient) getBoot(ctx context.Context, serverID int) (*BootProfile, error) {
	bytes, err := c.makeAPICall(ctx, "GET", fmt.Sprintf("%s/boot/%d", c.url, serverID), nil, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		return nil, err
	}

	jsonStr := string(bytes)
	bootProfile := BootProfile{}
	activeBoot := ""

	if gjson.Get(jsonStr, "boot.linux.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.linux").String()
		bootProfile.ActiveProfile = "linux"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.rescue.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.rescue").String()
		bootProfile.ActiveProfile = "rescue"
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}

	bootProfile.Architecture = gjson.Get(activeBoot, "arch").String()
	// bootProfile.AuthorizedKeys = gjson.Get(activeBoot, "authorised_keys").Array()
	// bootProfile.HostKeys = gjson.Get(activeBoot, "host_keys").Array()
	bootProfile.Password = gjson.Get(activeBoot, "password").String()
	bootProfile.ServerID = int(gjson.Get(activeBoot, "server_num").Int())
	bootProfile.ServerIPv4 = gjson.Get(activeBoot, "server_ip").String()
	bootProfile.ServerIPv6 = gjson.Get(activeBoot, "server_ipv6_net").String()

	return &bootProfile, nil
}

func (c *HetznerRobotClient) setBootProfile(ctx context.Context, serverID int, activeBootProfile string, arch string, os string, lang string, authorizedKeys []string) (*BootProfile, error) {
	data := url.Values{}
	data.Set("arch", arch)
	for _, key := range authorizedKeys {
		data.Add("authorized_key", key)
	}
	if activeBootProfile == "linux" {
		data.Set("dist", os)
		data.Set("lang", lang)
	}
	if activeBootProfile == "rescue" {
		data.Set("os", os)
	}

	bytes, err := c.makeAPICall(ctx, "POST", fmt.Sprintf("%s/boot/%d/%s", c.url, serverID, activeBootProfile), data, []int{http.StatusOK, http.StatusAccepted})
	if err != nil {
		if strings.Contains(err.Error(), "BOOT_ALREADY_ENABLED") {
			return c.getBoot(ctx, serverID)
		}
		return nil, err
	}

	jsonStr := string(bytes)
	bootProfile := BootProfile{}
	activeBoot := ""

	if gjson.Get(jsonStr, "boot.linux.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.linux").String()
		bootProfile.ActiveProfile = "linux"
		bootProfile.Language = gjson.Get(activeBoot, "lang").String()
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "dist").String()
	}
	if gjson.Get(jsonStr, "boot.rescue.active").Bool() {
		activeBoot = gjson.Get(jsonStr, "boot.rescue").String()
		bootProfile.ActiveProfile = "rescue"
		bootProfile.OperatingSystem = gjson.Get(activeBoot, "os").String()
	}

	bootProfile.Architecture = gjson.Get(activeBoot, "arch").String()
	// bootProfile.AuthorizedKeys = gjson.Get(activeBoot, "authorised_keys").Array()
	// bootProfile.HostKeys = gjson.Get(activeBoot, "host_keys").Array()
	bootProfile.Password = gjson.Get(activeBoot, "password").String()
	bootProfile.ServerID = int(gjson.Get(activeBoot, "server_num").Int())
	bootProfile.ServerIPv4 = gjson.Get(activeBoot, "server_ip").String()
	bootProfile.ServerIPv6 = gjson.Get(activeBoot, "server_ipv6_net").String()

	return &bootProfile, nil
}
