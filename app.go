package app

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	version "./version"
)

type PID int
type Path string

type ApplicationData struct {
	Local  Path
	Config Path
}

type Process struct {
	ID       PID
	Data     ApplicationData
	IO       IO
	Children map[PID]*Process
	Signals  chan *signal.Signal
}

type IO struct {
	Output *io.Writer
	Error  *io.Writer
	Input  *io.Reader
}

type Application struct {
	Name    string
	Version version.Version
	Process *Process
}

func New() *Application {
	fmt.Println("create the app object to be used by the service or the cli")
	fmt.Println("tool, or whatever. This is where the core application logic")
	fmt.Println("goes.\n")
	fmt.Println("this should be the location of code starting the process")
	fmt.Println("creating the PID file (could be library), supervision of ")
	fmt.Println("child processes (could be librar), signals (could be library)")
	fmt.Println("the very least it should be where these parts come together\n")
	fmt.Println("server logic likely should be in a library, but it depends")
	fmt.Println("on the application requirements, and service it provides\n")
	fmt.Println("config should be in its own library, ")
	return &Application{
		Name:    "app",
		Version: version.Semantic(0, 1, 0),
		IO: &IO{
			Output: os.Stdout(),
			Input:  os.Stdin(),
			Error:  os.Stderr(),
		},
	}
}
