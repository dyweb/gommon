package color

import (
	"fmt"
	"testing"
)

func Test_Colors(t *testing.T) {
	t.Skip("only used for trying out color codes")

	//fmt.Printf("\x1b[%dm%s\x1b[0m\n", nocolor, "nocolor")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", RedCode, "error")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", GreenCode, "ok")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", YellowCode, "warn")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", BlueCode, "info")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", PurpleCode, "京紫")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", CyanCode, "k1")
	fmt.Printf("\x1b[%dm%s\x1b[0m\n", GrayCode, "gray")
}

func TestRed(t *testing.T) {
	fmt.Println(Red("red"), "cry me a river", Cyan("no"))
	fmt.Println(Yellow("call me"), "maybe", Green("he he"))
}
