package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

var HomeDir string
var KasiopeaHomeDir string

func init() {
	var err error
	HomeDir, err = HomePath()
	if err != nil {
		ioutil.WriteFile("error.log", []byte(err.Error()), 0664)
		panic(err)
	} else {
		HomeDir += string(os.PathSeparator)
	}
	KasiopeaHomeDir = filepath.Join(HomeDir, "Documents", "kasiopea")
}
