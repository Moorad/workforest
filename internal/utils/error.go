package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Panic(msg string, err error) {
	color.Red(msg+":\n%s", err)
	os.Exit(1)
}

func PanicMsg(msg string, args ...any) {
	color.Red(msg, args...)
	os.Exit(1)
}

func GracefullyExit(msg string) {
	fmt.Println(msg)
	os.Exit(0)
}
