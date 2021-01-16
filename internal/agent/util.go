package agent

import (
	"gopkg.in/ini.v1"
)

type osRelease struct {
	ID string
}

func readOSRelease() (*osRelease, error) {
	osReleaseFile := "/etc/os-release"

	cfg, err := ini.Load(osReleaseFile)
	if err != nil {
		return nil, err
	}

	info := &osRelease{}
	info.ID = cfg.Section("").Key("ID").String()

	return info, nil
}
