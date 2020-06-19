package process

type PID int

type Process struct {
	ID       PID
	Data     ApplicationData
	IO       IO
	Children map[PID]*Process
	Signals  chan *signal.Signal
}
