package api

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	conf "github.com/sipt/kasiopea/config"
	. "github.com/sipt/kasiopea/constant"
	"github.com/sipt/kasiopea/extension/config"
	"github.com/sipt/kasiopea/upgrade"
)

var latest string
var url string
var status string

func CheckUpdate(ctx *gin.Context) {
	var err error
	latest, url, status, err = upgrade.CheckUpgrade(conf.KasiopeaVersion)
	if err != nil {
		ctx.JSON(500, Response{
			Code: 1, Message: err.Error(),
		})
		return
	}
	ctx.JSON(200, Response{
		Code: 0,
		Data: map[string]string{
			"Current": conf.KasiopeaVersion,
			"Latest":  latest,
			"URL":     url,
			"Status":  status,
		},
	})
}

func NewUpgrade(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		if status == upgrade.VersionEqual || status == upgrade.VersionGreater {
			ctx.JSON(500, Response{
				Code: 1, Message: "You're up-to-date!",
			})
			return
		}
		path := filepath.Join(config.HomeDir, "Downloads", "kasiopea.zip")
		os.Remove(path)
		err := upgrade.DownloadFile(path, url)
		if err != nil {
			ctx.JSON(500, Response{
				Code: 1, Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, Response{
			Code: 0, Message: "success",
		})
		eventChan <- EventUpgrade.SetData("kasiopea.zip")
	}
}
