package main

import (
	"fmt"
	"log"
	"os/exec"
)

// sudo grep "MARKB" ./demo.txt | tr -dc "0-9\n"
// sudo grep "Name=\"MARKB" /proc/bus/input/devices | cut -d '"' -f 2
// sudo awk '/MARKB-[0-9]+ Keyboard/ {p=1} NF==0 {p=0}; p' /proc/bus/input/devices

func GetHandlers() {
	out, err := exec.Command("/bin/sh", "-c", "awk '/MARKB-[0-9]+ Keyboard/ {p=1} NF==0 {p=0}; p' /proc/bus/input/devices").Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The date is %s\n", out)
}
