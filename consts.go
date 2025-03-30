package main

type AppState int
type TimerState int

//go:generate stringer -type=Programs
type Programs int

const (
	Open AppState = iota
	Closed
)

const (
	anki Programs = iota
	mpv
	ttsu
	asbplayer
	VLC
)

const (
	Play TimerState = iota
	Pause
	Running
)
