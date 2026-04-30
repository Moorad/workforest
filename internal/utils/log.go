package utils

import (
	"os"

	"github.com/fatih/color"
)

func Debug(msg string, args ...any) {
	debugEnv, _ := os.LookupEnv("DEBUG")

	if debugEnv == "1" {
		gray := color.RGB(170, 170, 170)

		_, err := gray.Printf(msg, args...)
		if err != nil {
			Panic("Failed to print debug log", err)
		}
		println()
	}
}

func DebugSuccess(msg string, args ...any) {
	debugEnv, _ := os.LookupEnv("DEBUG")

	if debugEnv == "1" {
		gray := color.RGB(103, 139, 100)

		_, err := gray.Printf(msg, args...)
		if err != nil {
			Panic("Failed to print debug log", err)
		}
		println()
	}
}

func DebugError(msg string, args ...any) {
	debugEnv, _ := os.LookupEnv("DEBUG")

	if debugEnv == "1" {
		gray := color.RGB(190, 142, 142)

		_, err := gray.Printf(msg, args...)
		if err != nil {
			Panic("Failed to print debug log", err)
		}
		println()
	}
}
