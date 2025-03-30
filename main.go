package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	dn "github.com/mitchellh/go-ps"
	"github.com/sevlyar/go-daemon"
	"log"
	"os/exec"
	"strings"
	"time"
)

type AnkiTimer struct {
	time  time.Duration
	state TimerState
	start time.Time
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

	timerch := make(chan TimerState)
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatal("DB open error: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB connection failed: ", err)
	}

	log.Printf("pasted opening")
	//initDB(db)
	log.Printf("pasted initDB")
	timer := AnkiTimer{
		time:  0,
		state: Pause,
	}

	for {
		if isAnkiRunning() {
			if timer.state == Pause {
				timer.start = time.Now()
				go trackAnkiTime(&timer, timerch)
			}
		} else {
			if timer.state == Running {
				timerch <- Pause
			}
		}
		if timer.state == Running {
			log.Printf("Time on Anki: %v", timer.time+time.Since(timer.start))
		} else {
			log.Printf("Time on Anki: %v", timer.time)
		}

		//sendDatatoDB(db, Anki, &timer)
		time.Sleep(500 * time.Millisecond)
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
	log.Printf(str)

	if str == s {
		return true
	}

	return false
}

// go routine for starting timer
func trackAnkiTime(timer *AnkiTimer, c <-chan TimerState) {
	start := timer.start

	timer.state = Running
	select {
	case msg := <-c:
		if msg == Pause {
			timer.time += time.Since(start)
			timer.state = Pause
			return
		}
	}
}

func initDB(db *sql.DB) {
	createProgramTableSQL := `
CREATE TABLE IF NOT EXISTS programs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT UNIQUE NOT NULL
);`

	createSessionTableSQL := `
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	program_id INTEGER NOT NULL REFERENCES programs(id),
	date DATE NOT NULL,
	duration INTEGER NOT NULL,
	UNIQUE(program_id, date)
);`

	_, err := db.Exec(createProgramTableSQL)
	if err != nil {
		log.Fatal("Error creating programs table:", err)
	}

	_, err = db.Exec(createSessionTableSQL)
	if err != nil {
		log.Fatal("Error creating sessions table:", err)
	}

	log.Println("Tables initialized successfully")
}

func sendDatatoDB(db *sql.DB, p Programs, t *AnkiTimer) {
	res, err := db.Exec(`INSERT OR IGNORE INTO programs (name) VALUES (?)`, "%s", p)
	if err != nil {
		log.Fatal(err)
	}
	programID, err := res.LastInsertId()

	_, err = db.Exec(`
INSERT INTO sessions (program_id, date, duration)
VALUES (?, ?, ?)
ON CONFLICT(program_id, date) DO UPDATE SET
    duration = excluded.duration
`, programID, time.Now().Format(time.DateOnly), t.time.String())

	if err != nil {
		log.Fatal(err)
	}
}
