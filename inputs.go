package main

import (
	"fmt"
	"log"
	"os/exec"
)

type HandlerStruct struct {
	name string
}

func GetHandlersRaw() string {
	out, err := exec.Command("/bin/sh", "-c", "awk '/MARKB-[0-9]+ Keyboard/ {p=1} NF==0 {p=0}; p' /proc/bus/input/devices").Output()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s\n", out)
}

func ParseHandlers() []HandlerStruct {
	var result []HandlerStruct

	return result
}
