package main

import (
	"fmt"
	dn "github.com/mitchellh/go-ps"
	"os/exec"
	"strings"
	"time"
)

func main() {

	time.Sleep(5 * time.Second)

	allPro, _ := dn.Processes()
	isAnkirunning := false
	isAnkiActive := false

	for _, pro := range allPro {

		if pro.Executable() == "anki" {
			isAnkirunning = true
		}
	}

	if isWindowActive("anki") {
		isAnkiActive = true
	} else {
		println("Anki not active lel")
	}

	if isAnkirunning && isAnkiActive {
		for {
			start := time.Now()

			time.Sleep(2 * time.Second)

			elapsed := time.Since(start)
			fmt.Printf("Elapsed time: %s\n", elapsed)

			if !isWindowActive("anki") {
				println("exited Anki stopping timer!")
				break
			}

		}
	} else {
		fmt.Printf("------------CANT START TIMER------------\n")
		fmt.Printf("Anki is running: %v \nAnki is active: %v", isAnkirunning, isAnkiActive)
	}

}

func isWindowActive(s string) bool {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to get name of first application process whose frontmost is true`)
	out, _ := cmd.Output()
	str := string(out)
	str = strings.TrimSpace(str)

	if str == s {
		return true
	}

	return false
}
