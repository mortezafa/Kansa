package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	dn "github.com/mitchellh/go-ps"
	"github.com/sevlyar/go-daemon"
)

type ProgramTimer struct {
	time    time.Duration
	state   TimerState
	start   time.Time
	program Programs
}

func main() {

	if len(os.Args) > 1 && os.Args[1] == "report" {
		runReport()
		return
	}

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

	db, err := sql.Open("sqlite3", "./kansa.db")

	if err != nil {
		log.Fatal("DB open error: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB connection failed: ", err)
	}
	defer db.Close()

	log.Printf("pasted opening")
	initDB(db)
	log.Printf("pasted initDB")
	timers := map[Programs]*ProgramTimer{
		anki:      &ProgramTimer{program: anki, state: Pause, start: time.Now()},
		mpv:       &ProgramTimer{program: mpv, state: Pause, start: time.Now()},
		ttsu:      &ProgramTimer{program: ttsu, state: Pause, start: time.Now()},
		asbplayer: &ProgramTimer{program: asbplayer, state: Pause, start: time.Now()},
		VLC:       &ProgramTimer{program: VLC, state: Pause, start: time.Now()},
	}

	for {
		currentProg := getCurrentProg()

		for _, prog := range timers {
			if isProgRunning(prog.program.String()) && prog.program.String() == currentProg && prog.state == Pause {
				prog.start = time.Now()
				prog.state = Running
			} else {
				if prog.state == Running {
					prog.time += time.Since(prog.start)
					prog.state = Pause
					sendDatatoDB(db, prog.program, prog)
				}
			}
			if prog.state == Running {
				log.Printf("Time on %s: %v \n", prog.program.String(), prog.time+time.Since(prog.start))
			} else {
				log.Printf("Time on %s: %v \n", prog.program.String(), prog.time)
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

}

func isProgRunning(pn string) bool {
	allPro, _ := dn.Processes()

	for _, pro := range allPro {
		if pro.Executable() == pn && isWindowActive(pn) {
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
	log.Printf("CURRENT PROG: %s", str)

	if str == s {
		return true
	}

	return false
}

func getCurrentProg() string {
	cmd := exec.Command("osascript", "-e", `tell application "System Events" to get name of first application process whose frontmost is true`)
	out, _ := cmd.Output()
	str := string(out)
	str = strings.TrimSpace(str)
	return str
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

func sendDatatoDB(db *sql.DB, p Programs, t *ProgramTimer) {
	res, err := db.Exec(`INSERT OR IGNORE INTO programs (name) VALUES (?)`, p.String())
	if err != nil {
		log.Fatal(err)
	}
	programID, err := res.LastInsertId()
	if programID == 0 {
		// row already existed, look it up
		row := db.QueryRow(`SELECT id FROM programs WHERE name = ?`, p.String())
		if err := row.Scan(&programID); err != nil {
			log.Fatal(err)
		}
	}

	_, err = db.Exec(`
INSERT INTO sessions (program_id, date, duration)
VALUES (?, ?, ?)
ON CONFLICT(program_id, date) DO UPDATE SET
    duration = excluded.duration
`, programID, time.Now().Format(time.DateOnly), int64(t.time.Seconds()))

	if err != nil {
		log.Fatal(err)
	}
}

func runReport() {
	db, err := sql.Open("sqlite3", "./kansa.db")
	if err != nil {
		log.Fatal("DB open error: ", err)
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT p.name, s.date, s.duration
FROM sessions s
JOIN programs p ON p.id = s.program_id
ORDER BY s.date DESC, p.name ASC
`)
	if err != nil {
		log.Fatal("Query error: ", err)
	}
	defer rows.Close()

	fmt.Println("Kansa immersion report")
	fmt.Println("-----------------------")
	for rows.Next() {
		var name, date string
		var seconds int64
		if err := rows.Scan(&name, &date, &seconds); err != nil {
			log.Fatal(err)
		}
		d := time.Duration(seconds) * time.Second
		fmt.Printf("%-10s %-12s %v\n", name, date, d)
	}
}
