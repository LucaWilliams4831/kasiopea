package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/sipt/kasiopea/config"
	. "github.com/sipt/kasiopea/constant"
	"github.com/sipt/kasiopea/controller"
	"github.com/sipt/kasiopea/log"
)

var eventChan chan *EventObj

func ListenEvent() {
	eventChan = make(chan *EventObj, 1)
	go dealEvent(eventChan)
}

func dealEvent(c chan *EventObj) {
	for {
		t := <-c
		switch t.Type {
		case EventShutdown.Type:
			log.Logger.Info("[Kasiopea] is shutdown, see you later!")
			shutdown(config.CurrentConfig().General.SetAsSystemProxy)
			os.Exit(0)
			return
		case EventReloadConfig.Type:
			_, err := reloadConfig(config.CurrentConfigFile(), StopSocksSignal, StopHTTPSignal)
			if err != nil {
				log.Logger.Error("Reload Config failed: ", err)
				fmt.Println(err.Error())
				os.Exit(1)
			}
		case EventRestartHttpProxy.Type:
			StopHTTPSignal <- true
			go HandleHTTP(config.CurrentConfig(), StopHTTPSignal)
		case EventRestartSocksProxy.Type:
			StopSocksSignal <- true
			go HandleSocks5(config.CurrentConfig(), StopSocksSignal)
		case EventRestartController.Type:
			controller.ShutdownController()
			go controller.StartController(config.CurrentConfig(), eventChan)
		case EventUpgrade.Type:
			//todo
			fileName := t.GetData().(string)
			shutdown(config.CurrentConfig().General.SetAsSystemProxy)
			log.Logger.Info("[Kasiopea] is shutdown, for upgrade!")
			var name string
			if runtime.GOOS == "windows" {
				name = "upgrade"
			} else {
				name = "./upgrade"
			}
			cmd := exec.Command(name, "-f="+fileName)
			err := cmd.Start()
			if err != nil {
				fmt.Println(err.Error())
			}
			os.Exit(0)
		}
	}
}
