package goadb

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_ADB_PATH = "/xinxi/adb"
)

func ForTestCommandExecuter(cmd string, args ...string) (string, error) {
	fullCmd := append([]string{cmd}, args...)
	return strings.Join(fullCmd, " "), nil
}

func TestGeneralCommand(t *testing.T) {
	device1 := &Device{
		transportID:     1,
		adbPath:         TEST_ADB_PATH,
		commandExecuter: ForTestCommandExecuter,
	}
	device2 := &Device{
		transportID:     2,
		adbPath:         TEST_ADB_PATH,
		commandExecuter: ForTestCommandExecuter,
	}

	testTable := []struct {
		testIndex   int
		device      *Device
		inputCmd    string
		expectedCmd string
	}{
		{1, device1, "whoami", TEST_ADB_PATH + " -t 1 shell whoami"},
		{2, device1, "/axon/bin/axctl dvr -state", TEST_ADB_PATH + " -t 1 shell /axon/bin/axctl dvr -state"},
		{3, device2, "whoami", TEST_ADB_PATH + " -t 2 shell whoami"},
		{4, device2, "/axon/bin/axctl dvr -state", TEST_ADB_PATH + " -t 2 shell /axon/bin/axctl dvr -state"},
	}

	for _, testCase := range testTable {
		output, err := testCase.device.ShellCmd(testCase.inputCmd)
		assert.Nil(t, err)
		assert.Equal(t, string(output), testCase.expectedCmd)
	}
}

func TestGetDevice(t *testing.T) {
	adbPath := findAdbPath()
	if len(adbPath) == 0 {
		return // this env does not have adb installed, abort
	}
	devices, err := getDevices(adbPath)
	assert.Nil(t, err)
	if len(devices) == 0 {
		return // current env does not have adb device connected
	}
	for i, device := range devices {
		assert.Equal(t, device.transportID, i+1)
		hostname, err := device.HostName()
		assert.Nil(t, err)
		assert.NotEqual(t, len(hostname), 0)
	}
}
