package application

import (
	"./filesystem"
)

type ApplicationData struct {
	Local  filesystem.Path
	Config filesystem.Path
}
