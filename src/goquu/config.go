package goquu
import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os/user"
	"path"
	"errors"
	"log"
)

type Config struct {
	LogFile string
	LogFlag int
	WorkerSize int
	DBDirectory string
}

func (config *Config) setDefault() {
	config.LogFile = "/tmp/goquu.log"
	config.DBDirectory = "/tmp/db"
	config.WorkerSize = 1
	config.LogFlag = log.LstdFlags
}

func load(path string) (config *Config, err error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	config = &Config{}
	config.setDefault()
	result := *config
	_, err = toml.Decode(string(text), &result)

	if err != nil {
		return
	}
	return &result, err
}

func loadConfig() (config *Config, err error) {
	usr, err := user.Current()
	if err == nil {
		config, err = load(path.Join(usr.HomeDir, ".goquu.toml"))
	}
	if err == nil {
		return
	}
	config, err = load("/etc/goquu.toml")
	if err == nil {
		return
	}
	return config, errors.New("Cannot find any readable configuration files!")
}

func (server *Server) OnClose() {
}
