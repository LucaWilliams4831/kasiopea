package conf

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/kasiopea"
)

func GetMitMRules(ctx *gin.Context) {
	var response Response
	response.Data = kasiopea.GetMitMRules()
	ctx.JSON(200, response)
}

func AppendMitMRules(ctx *gin.Context) {
	d := ctx.Query("domain")
	if len(d) > 0 {
		kasiopea.AppendMitMRules(d)
	}
	var response Response
	response.Data = kasiopea.GetMitMRules()
	ctx.JSON(200, response)
}

func DelMitMRules(ctx *gin.Context) {
	d := ctx.Query("domain")
	if len(d) > 0 {
		kasiopea.RemoveMitMRules(d)
	}
	var response Response
	response.Data = kasiopea.GetMitMRules()
	ctx.JSON(200, response)
}
