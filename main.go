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

	for i := range inputDevices {
		fmt.Printf("Name=%s\n", inputDevices[i].name)
		fmt.Printf("Phys=%s\n", inputDevices[i].phys)
		for j := range inputDevices[i].handlers {
			fmt.Printf("Handlers[%d]=%s\n", j, inputDevices[i].handlers[j])
		}
		fmt.Println("///   ///   ///")
	}

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
	}

	if len(inputDevices) > 0 {
		for i := range inputDevices[0].handlers {
			if strings.ContainsAny(inputDevices[0].handlers[i], "event") {
				readLoop(inputDevices[0].handlers[i], inputDevices[0].phys)
			}
		}
	}

	if token := c.Unsubscribe("go-mqtt/sample"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
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
				fmt.Printf("%s %s\n", physicalAddress, convertSequenceToString(input))
				input = []uint16{}
			}
		}

		//fmt.Printf("type: %x\ncode: %d\nvalue: %d\n", typ, code, value)
	}
}

func convertSequenceToString(keycodeSequence []uint16) string {
	var result string = ""

	var uppercase bool = false
	for i := range keycodeSequence {
		var char string = ""

		if keycodeSequence[i] == KEY_LEFTSHIFT {
			uppercase = true
			continue
		}

		switch keycodeSequence[i] {
		case KEY_A:
			char = "a"
		case KEY_B:
			char = "b"
		case KEY_C:
			char = "c"
		case KEY_D:
			char = "d"
		case KEY_E:
			char = "e"
		case KEY_F:
			char = "f"
		case KEY_G:
			char = "g"
		case KEY_H:
			char = "h"
		case KEY_I:
			char = "i"
		case KEY_J:
			char = "j"
		case KEY_K:
			char = "k"
		case KEY_L:
			char = "l"
		case KEY_M:
			char = "m"
		case KEY_N:
			char = "n"
		case KEY_O:
			char = "o"
		case KEY_P:
			char = "p"
		case KEY_Q:
			char = "q"
		case KEY_R:
			char = "r"
		case KEY_S:
			char = "s"
		case KEY_T:
			char = "t"
		case KEY_U:
			char = "u"
		case KEY_V:
			char = "v"
		case KEY_W:
			char = "w"
		case KEY_X:
			char = "x"
		case KEY_Y:
			char = "y"
		case KEY_Z:
			char = "z"
		case KEY_1:
			char = "1"
		case KEY_2:
			char = "2"
		case KEY_3:
			char = "3"
		case KEY_4:
			char = "4"
		case KEY_5:
			char = "5"
		case KEY_6:
			char = "6"
		case KEY_7:
			char = "7"
		case KEY_8:
			char = "8"
		case KEY_9:
			char = "9"
		case KEY_0:
			char = "0"
		case KEY_MINUS:
			char = "-"
		}

		if uppercase == true {
			char = strings.ToUpper(char)
			uppercase = false
		}

		result += char
	}

	return result
}
