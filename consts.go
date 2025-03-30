package main

type AppState int
type TimerState int

type Programs int

const (
	Open AppState = iota
	Closed
)

const (
	Anki Programs = iota
)

const (
	Play TimerState = iota
	Pause
	Running
)
