package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	evdev "github.com/gvalkov/golang-evdev"
	"os"
	"strings"
	"syscall"
)

const BrokerAddress = "tcp://leebapp1.leeb.cc:1883"
const BrokerChannel = "leeb/gnesau/verladung/scanner"

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	inputDevices := GetInputDevices()

	fmt.Println("###   Registered Devices   ###")
	for i := range inputDevices {
		fmt.Printf("%s (%s)\n", inputDevices[i].name, inputDevices[i].uniq)
	}
	fmt.Println("###                        ###")

	opts := MQTT.NewClientOptions().AddBroker(BrokerAddress)
	opts.SetClientID("ScanHanGo_" + hostname)
	opts.SetDefaultPublishHandler(f)

	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	inputDevices = []InputDevice{}

	var currentLoops []string

	for {
		newInputDevices := GetInputDevices()

		if len(inputDevices) != len(newInputDevices) {
			inputDevices = newInputDevices
			for i := range inputDevices {
				if !contains(currentLoops, inputDevices[i].uniq) {
					for j := range inputDevices[i].handlers {
						if strings.Contains(inputDevices[i].handlers[j], "event") {
							fmt.Printf("NEW DEVICE READING: %s (%s)\n", inputDevices[i].name, inputDevices[i].uniq)
							go readLoop(inputDevices[i].handlers[j], inputDevices[i].uniq, c)
						}
					}
				}
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func readLoop(event string, physicalAddress string, mqttClient MQTT.Client) {
	var input []uint16
	f, err := os.Open("/dev/input/" + event)
	if err != nil {
		panic(err)
	}
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), evdev.EVIOCGRAB, 1)
	defer f.Close()
	defer syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), evdev.EVIOCGRAB, 0)
	b := make([]byte, 24)
	for {
		f.Read(b)
		var value int32
		typ := binary.LittleEndian.Uint16(b[16:18])
		code := binary.LittleEndian.Uint16(b[18:20])
		binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)

		if typ == 1 && value == 0 {
			input = append(input, code)
		}

		if len(input) > 0 {
			if input[len(input)-1] == 28 {
				text := ConvertSequenceToString(input)
				fmt.Printf("%s %s\n", physicalAddress, text)
				mqttClient.Publish(BrokerChannel+"/"+physicalAddress+"/scanned", 0, false, text)
				input = []uint16{}
			}
		}
	}
}
