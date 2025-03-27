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

	cmd := exec.Command("osascript", "-e", `tell application "System Events" to get name of first application process whose frontmost is true`)
	out, _ := cmd.Output()
	s := string(out)
	s = strings.TrimSpace(s)

	if s == "anki" {
		isAnkiActive = true
	} else {
		println(string(out))
	}

	if isAnkirunning && isAnkiActive {
		for {
			start := time.Now()

			time.Sleep(2 * time.Second)

			elapsed := time.Since(start)
			fmt.Printf("Elapsed time: %s\n", elapsed)
		}
	} else {
		fmt.Printf("------------CANT START TIMER------------\n")
		fmt.Printf("Anki is running: %v \nAnki is active: %v", isAnkirunning, isAnkiActive)
	}

}
