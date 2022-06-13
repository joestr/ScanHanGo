package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

type InputDevice struct {
	name     string
	phys     string
	handlers []string
}

func GetInputDevices() []InputDevice {
	var result []InputDevice

	var singleDevices []string = SplitIntoSingleDevicesRaw(GetInputDevicesRaw())

	for i := range singleDevices {
		result = append(result, ParseInputDevice(singleDevices[i]))
	}

	return result
}

func ParseInputDevice(device string) InputDevice {
	var result InputDevice

	regex := regexp.MustCompile("\n")
	deviceLines := regex.Split(device, -1)

	isNameSet := false
	isPhysSet := false
	areHandlersSet := false

	for i := range deviceLines {
		if !isNameSet {
			val, success := ParseInputDeviceName(deviceLines[i])
			if success {
				result.name = val
				isNameSet = true
				continue
			}
		}
		if !isPhysSet {
			val, success := ParseInputDevicePhys(deviceLines[i])
			if success {
				result.phys = val
				isPhysSet = true
				continue
			}
		}
		if !areHandlersSet {
			val, success := ParseInputDeviceHandlers(deviceLines[i])
			if success {
				result.handlers = val
				areHandlersSet = true
				break
			}
		}
	}

	return result
}

func GetInputDevicesRaw() string {
	out, err := exec.Command("/bin/sh", "-c", "awk '/MARKB-[0-9]+ Keyboard/ {p=1} NF==0 {p=0}; p' /proc/bus/input/devices").Output()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s\n", out)
}

func SplitIntoSingleDevicesRaw(devices string) []string {
	regex := regexp.MustCompile("^\n$")

	split := regex.Split(devices, -1)
	var set []string

	for i := range split {
		set = append(set, split[i])
	}

	return set
}

func ParseInputDeviceName(line string) (string, bool) {
	var result = ""

	parts := strings.Split(line, "Name=")
	if len(parts) < 2 {
		return result, false
	}

	result = strings.Replace(parts[1], "\"", "", -1)

	return result, true
}

func ParseInputDevicePhys(line string) (string, bool) {
	var result = ""

	parts := strings.Split(line, "Phys=")
	if len(parts) < 2 {
		return result, false
	}

	result = parts[1]

	return result, true
}

func ParseInputDeviceHandlers(line string) ([]string, bool) {
	var result []string

	parts := strings.Split(line, "Handlers=")
	if len(parts) < 2 {
		return result, false
	}

	result = strings.Split(parts[1], " ")

	return result, true
}
