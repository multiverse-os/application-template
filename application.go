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
	app := &Application{
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

	if _, err := os.Stat(app.Data.Path); os.IsNotExist(err) {
		_ = os.MkdirAll(app.Data.Path, os.FileMode(0770))
	}

	if _, err := os.Stat(app.Config.Path); os.IsNotExist(err) {
		_ = os.MkdirAll(app.Config.Path, os.FileMode(0770))
	}

	return app
}
