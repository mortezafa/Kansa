package main

import (
	dn "github.com/mitchellh/go-ps"
	"github.com/sevlyar/go-daemon"
	"log"
	"os/exec"
	"strings"
	"time"
)

type AppState int
type TimerState int

const (
	Open AppState = iota
	Closed
)

const (
	Play TimerState = iota
	Pause
	Running
)

type AnkiTimer struct {
	time  time.Duration
	state TimerState
}

func main() {

	cntxt := &daemon.Context{
		PidFileName: "kansa.pid",
		PidFilePerm: 0644,
		LogFileName: "kansa.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()

	log.Printf("Daemon Started!!!! Kansa daemon")

	for {
		log.Print("In main loop")

		// idle state

		// anki open + actve

		// anki closed (stopped timer)

	}

}

// go routine for wathcing anki
func isAnkiRunning() bool {
	allPro, _ := dn.Processes()

	for _, pro := range allPro {

		if pro.Executable() == "anki" && isWindowActive("anki") {
			return true
		}
	}

	return false

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

// go routine for starting timer
func trackAnkiTime(timer *AnkiTimer) {
	if timer.state == Running {
		var elapsed time.Duration
		start := time.Now()
		timer.state = Running

		elapsed += time.Since(start)
		timer.time += elapsed
	}
}
