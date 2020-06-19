package application

import (
	"fmt"
	"os"

	"./filesystem"
)

type Application struct {
	Name    string
	Version Version
	Process *Process
	IO      *IO
	Data    filesystem.Directory
	Config  filesystem.Directory
}

func Initialize(name string, version Version) *Application {
	if len(name) == 0 {
		name = os.Executable()
	}

	return &Application{
		Name:    name,
		Version: version,
		IO: &IO{
			Output: os.Stdout(),
			Input:  os.Stdin(),
			Error:  os.Stderr(),
		},
		Data: filesystem.Directory{
			Path: filesystem.Path{fmt.Sprintf("/home/%s/%s/%s", os.Getenv("USER"), ".local/share", name)},
		},
		Config: filesystem.Directory{
			Path: filesystem.Path{fmt.Sprintf("/home/%s/%s/%s", os.Getenv("USER"), ".config", name)},
		},
	}
}
