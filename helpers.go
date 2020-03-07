package goadb

import (
	"fmt"
	"os/exec"
	"strings"
)

func findAdbPath() string {
	cmd := exec.Command("which", "adb")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(output[:len(output)-1]) // get rid of \n at the end
}

func getDevices(adbPath string) (devices []*Device, err error) {
	cmd := exec.Command(adbPath, "start-server")
	_, err = cmd.CombinedOutput()
	if err != nil {
		return
	}
	cmd = exec.Command(adbPath, "devices")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if len(line) > 0 {
			if strings.Contains(line, "List") {
				continue
			}
			deviceLine := strings.Split(line, "	")
			if len(deviceLine) != 2 {
				return devices, fmt.Errorf("get-devices: cannot parse %s", line)
			}
			device := &Device{
				transportID: i,
				adbPath:     adbPath,
			}
			device.commandExecuter = device.executeShellCmd
			devices = append(devices, device)
		}
	}

	return
}
