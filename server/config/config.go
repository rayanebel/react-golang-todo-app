package config

import (
	"github.com/rayanebel/react-golang-todo-app/server/models"
)

var Config models.Config

func SaveConfig(config models.Config) {
	Config = config
	return
}
