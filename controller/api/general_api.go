package api

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sipt/kasiopea/config"
	. "github.com/sipt/kasiopea/constant"
	"github.com/sipt/kasiopea/extension/network"
	"github.com/sipt/kasiopea/rule"
)

func EnableSystemProxy(ctx *gin.Context) {
	g := config.CurrentConfig().General
	network.WebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SecureWebProxySwitch(true, "127.0.0.1", g.HttpPort)
	network.SocksProxySwitch(true, "127.0.0.1", g.SocksPort)
	ctx.JSON(200, Response{})
}

func DisableSystemProxy(ctx *gin.Context) {
	network.WebProxySwitch(false)
	network.SecureWebProxySwitch(false)
	network.SocksProxySwitch(false)
	ctx.JSON(200, Response{})
}

func NewShutdown(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		eventChan <- EventShutdown
	}
}

func ReloadConfig(eventChan chan *EventObj) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, Response{})
		eventChan <- EventReloadConfig
	}
}

func GetConnMode(ctx *gin.Context) {
	ctx.JSON(200, Response{
		Data: rule.GetConnMode(),
	})
}

func SetConnMode(ctx *gin.Context) {
	value := ctx.Param("mode")
	value = strings.ToUpper(value)
	err := rule.SetConnMode(value)
	if err != nil {
		ctx.JSON(500, Response{
			Code:    1,
			Message: err.Error(),
		})
	}
	GetConnMode(ctx)
}
