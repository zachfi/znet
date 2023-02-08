package agent

import (
	"io"
	"os/exec"

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

func runCommand(name string, arg ...string) (*CommandResult, error) {
	cmd := exec.Command(name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	outResult, _ := io.ReadAll(stdout)
	errResult, _ := io.ReadAll(stderr)

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	result := &CommandResult{
		Output:   outResult,
		Error:    errResult,
		ExitCode: int32(cmd.ProcessState.ExitCode()),
	}

	return result, nil
}
