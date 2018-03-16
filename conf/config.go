package conf

import (
	"encoding/json"
	"io/ioutil"
)

var (
	Conf       Servers
	configFile = "./conf/config.json"
)

type Servers struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Host         string    `json:"host"`
	Port         int64     `json:"port"`
	User         string    `json:"user"`
	Password     string    `json:"password"`
	PreCommands  string    `json:"preCommands"`
	Uploads      []Uploads `json:"uploads"`
	PostCommands string    `json:"postCommands"`
}

type Uploads struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

func InitConfig() (err error) {
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(config, &Conf)
	if err != nil {
		return err
	}
	return nil
}
