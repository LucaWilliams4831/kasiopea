package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sipt/kasiopea"
)

func GetRecords(ctx *gin.Context) {
	ctx.JSON(200, &Response{
		Data: kasiopea.GetRecords(),
	})
}
func ClearRecords(ctx *gin.Context) {
	kasiopea.ClearRecords()
	dump := kasiopea.GetDump()
	if dump != nil {
		dump.Clear()
	}
	ctx.JSON(200, &Response{})
}
