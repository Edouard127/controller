package controller

type Signal int

const (
	Reload = iota
	Shutdown
	Dump
	Debug
	Halt
	NewCircuit
	ClearCircuit
	Heartbeat
	Dormant
	Active
)

func (s Signal) String() string {
	return [...]string{
		"RELOAD",
		"SHUTDOWN",
		"DUMP",
		"DEBUG",
		"HALT",
		"NEWNYM",
		"CLEARDNSCACHE",
		"HEARTBEAT",
		"DORMANT",
		"ACTIVE",
	}[s]
}
