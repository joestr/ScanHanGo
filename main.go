package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"os"
	"strings"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {

	inputDevices := GetInputDevices()

	fmt.Println("###   Registered Devices   ###")
	for i := range inputDevices {
		fmt.Printf("%s (%s)\n", inputDevices[i].name, inputDevices[i].uniq)
	}
	fmt.Println("###                        ###")

	/*
		opts := MQTT.NewClientOptions().AddBroker("tcp://mqtt.eclipseprojects.io:1883")
		opts.SetClientID("go-simple")
		opts.SetDefaultPublishHandler(f)

		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := c.Subscribe("go-mqtt/sample", 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for i := 0; i < 5; i++ {
			text := fmt.Sprintf("this is msg #%d!", i)
			token := c.Publish("go-mqtt/sample", 0, false, text)
			token.Wait()
		}*/

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
							go readLoop(inputDevices[i].handlers[j], inputDevices[i].uniq)
						}
					}
				}
			}
		}
	}

	/*if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)*/
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func readLoop(event string, physicalAddress string) {
	var input []uint16
	f, err := os.Open("/dev/input/" + event)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b := make([]byte, 24)
	for {
		f.Read(b)
		//sec := binary.LittleEndian.Uint64(b[0:8])
		//usec := binary.LittleEndian.Uint64(b[8:16])
		//t := time.Unix(int64(sec), int64(usec))
		//fmt.Println(t)
		var value int32
		typ := binary.LittleEndian.Uint16(b[16:18])
		code := binary.LittleEndian.Uint16(b[18:20])
		binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &value)

		if typ == 1 && value == 0 {
			input = append(input, code)
		}

		if len(input) > 0 {
			if input[len(input)-1] == 28 {
				fmt.Printf("%s %s\n", physicalAddress, ConvertSequenceToString(input))
				input = []uint16{}
			}
		}

		//fmt.Printf("type: %x\ncode: %d\nvalue: %d\n", typ, code, value)
	}
}
