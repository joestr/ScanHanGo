package main

import (
	"github.com/holoplot/go-evdev"
	"strings"
)

func ConvertSequenceToString(keycodeSequence []uint16) string {
	var result string = ""

	var uppercase bool = false
	for i := range keycodeSequence {
		var char string = ""

		if keycodeSequence[i] == evdev.KEY_LEFTSHIFT {
			uppercase = true
			continue
		}

		switch keycodeSequence[i] {
		case evdev.KEY_A:
			char = "a"
		case evdev.KEY_B:
			char = "b"
		case evdev.KEY_C:
			char = "c"
		case evdev.KEY_D:
			char = "d"
		case evdev.KEY_E:
			char = "e"
		case evdev.KEY_F:
			char = "f"
		case evdev.KEY_G:
			char = "g"
		case evdev.KEY_H:
			char = "h"
		case evdev.KEY_I:
			char = "i"
		case evdev.KEY_J:
			char = "j"
		case evdev.KEY_K:
			char = "k"
		case evdev.KEY_L:
			char = "l"
		case evdev.KEY_M:
			char = "m"
		case evdev.KEY_N:
			char = "n"
		case evdev.KEY_O:
			char = "o"
		case evdev.KEY_P:
			char = "p"
		case evdev.KEY_Q:
			char = "q"
		case evdev.KEY_R:
			char = "r"
		case evdev.KEY_S:
			char = "s"
		case evdev.KEY_T:
			char = "t"
		case evdev.KEY_U:
			char = "u"
		case evdev.KEY_V:
			char = "v"
		case evdev.KEY_W:
			char = "w"
		case evdev.KEY_X:
			char = "x"
		case evdev.KEY_Y:
			char = "y"
		case evdev.KEY_Z:
			char = "z"
		case evdev.KEY_1:
			char = "1"
		case evdev.KEY_2:
			char = "2"
		case evdev.KEY_3:
			char = "3"
		case evdev.KEY_4:
			char = "4"
		case evdev.KEY_5:
			char = "5"
		case evdev.KEY_6:
			char = "6"
		case evdev.KEY_7:
			char = "7"
		case evdev.KEY_8:
			char = "8"
		case evdev.KEY_9:
			char = "9"
		case evdev.KEY_0:
			char = "0"
		case evdev.KEY_MINUS:
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
