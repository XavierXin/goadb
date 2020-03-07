package goadb

import (
	"errors"
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
	cmd = exec.Command(adbPath, "devices", "-l")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if len(line) > 0 {
			if strings.Contains(line, "List") {
				continue
			}
			deviceLine := strings.Split(line, "	")
			transportID := "0"
			for _, deviceLineSegment := range deviceLine {
				if strings.Contains(deviceLineSegment, "transport_id") {
					transportIDPair := strings.Split(deviceLineSegment, ":")
					if len(transportIDPair) >= 2 {
						transportID = transportIDPair[1]
						break
					}
				}
			}
			if transportID == "0" {
				return devices, errors.New("your adb version does not support connect by transport")
			}
			device := &Device{
				transportID: transportID,
				adbPath:     adbPath,
			}
			device.commandExecuter = device.executeShellCmd
			devices = append(devices, device)
		}
	}

	return
}
