package api

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/apaxa-go/helper/strconvh"
	"github.com/gin-gonic/gin"
	"github.com/sipt/kasiopea"
	"github.com/sipt/kasiopea/log"
)

func SetAllowDump(ctx *gin.Context) {
	var response Response
	allow_mitm := ctx.PostForm("allow_mitm")
	allow_dump := ctx.PostForm("allow_dump")
	switch allow_dump {
	case "true":
		kasiopea.SetAllowDump(true)
	case "false":
		kasiopea.SetAllowDump(false)
	case "":
	default:
		response.Code = 1
		response.Message = fmt.Sprintf("allow_dump value error: %v", allow_dump)
		ctx.JSON(500, response)
		return
	}
	switch allow_mitm {
	case "true":
		kasiopea.SetAllowMitm(true)
	case "false":
		kasiopea.SetAllowMitm(false)
	case "":
	default:
		response.Code = 1
		response.Message = fmt.Sprintf("allow_mitm value error: %v", allow_mitm)
		ctx.JSON(500, response)
		return
	}
	GetAllowDump(ctx)
}

func GetAllowDump(ctx *gin.Context) {
	var response = Response{
		Data: struct {
			AllowDump bool `json:"allow_dump"`
			AllowMitm bool `json:"allow_mitm"`
		}{
			kasiopea.GetAllowDump(),
			kasiopea.GetAllowMitm(),
		},
	}
	ctx.JSON(200, response)
}

func DumpRequest(ctx *gin.Context) {
	var response Response
	idStr := ctx.Param("conn_id")
	id, err := strconvh.ParseInt64(idStr)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	r := kasiopea.GetRecord(id)
	if r == nil {
		response.Code = 1
		response.Message = idStr + " not exist"
		ctx.JSON(500, response)
		return
	}
	if r.Status != kasiopea.RecordStatusCompleted {
		response.Code = 1
		response.Message = idStr + " not Completed"
		ctx.JSON(500, response)
		return
	}
	dump := kasiopea.GetDump()
	if dump == nil {
		response.Code = 1
		response.Message = "IDump is nil"
		ctx.JSON(500, response)
		return
	}
	reqStream, reqSize, respStream, respSize, err := dump.Dump(id)
	defer func() {
		if reqStream != nil {
			reqStream.Close()
		}
		if respStream != nil {
			respStream.Close()
		}
	}()
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	respStruct := &struct {
		ReqHeader  string
		ReqBody    string
		RespBody   string
		RespHeader string
	}{}
	buffer := &bytes.Buffer{}
	req, err := http.ReadRequest(bufio.NewReader(reqStream))
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	if reqSize > kasiopea.LargeRequestBody {
		buffer.WriteString("large body")
	} else {
		buffer.ReadFrom(req.Body)
		req.Body.Close()
	}
	respStruct.ReqBody = base64.StdEncoding.EncodeToString(buffer.Bytes())
	buffer.Reset()
	req.Write(buffer)
	respStruct.ReqHeader = base64.StdEncoding.EncodeToString(buffer.Bytes())
	buffer.Reset()

	resp, err := http.ReadResponse(bufio.NewReader(respStream), req)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}

	if respSize > kasiopea.LargeResponseBody {
		buffer.WriteString("large body")
		resp.Body.Close()
	} else {
		var r io.Reader
		if resp.Header.Get("Content-Encoding") == "gzip" {
			r, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Logger.Errorf("[Kasiopea-Controller] [%d] gzip init for response failed: %v", id, err)
				response.Code = 1
				response.Message = err.Error()
				ctx.JSON(500, response)
				return
			}
		} else if resp.Header.Get("Content-Encoding") == "deflate" {
			r, err = zlib.NewReader(resp.Body)
			if err != nil {
				log.Logger.Errorf("[Kasiopea-Controller] [%d] deflate init for response failed: %v", id, err)
				response.Code = 1
				response.Message = err.Error()
				ctx.JSON(500, response)
				return
			}
		} else {
			r = resp.Body
		}
		buffer.ReadFrom(r)
		resp.Body.Close()
	}
	respStruct.RespBody = base64.StdEncoding.EncodeToString(buffer.Bytes())
	buffer.Reset()
	resp.Write(buffer)
	respStruct.RespHeader = base64.StdEncoding.EncodeToString(buffer.Bytes())

	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.JSON(200, Response{
		Code:    0,
		Message: "success",
		Data:    respStruct,
	})
}

func DumpLarge(ctx *gin.Context) {
	response := Response{}
	fileName := ctx.Query("file_name")
	if len(fileName) == 0 {
		response.Code = 1
		response.Message = "file_name is empty!"
		ctx.JSON(500, response)
		return
	}

	dumpType := ctx.Query("dump_type")
	if len(dumpType) == 0 {
		response.Code = 1
		response.Message = "dump_type is empty!"
		ctx.JSON(500, response)
		return
	}

	if dumpType != "request" && dumpType != "response" {
		response.Code = 1
		response.Message = "dump_type must be 'request' or 'response'!"
		ctx.JSON(500, response)
		return
	}
	idStr := ctx.Param("conn_id")
	id, err := strconvh.ParseInt64(idStr)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	r := kasiopea.GetRecord(id)
	if r == nil {
		response.Code = 1
		response.Message = idStr + " not exist"
		ctx.JSON(500, response)
		return
	}
	if r.Status != kasiopea.RecordStatusCompleted {
		response.Code = 1
		response.Message = idStr + " not Completed"
		ctx.JSON(500, response)
		return
	}
	dump := kasiopea.GetDump()
	if dump == nil {
		response.Code = 1
		response.Message = "IDump is nil"
		ctx.JSON(500, response)
		return
	}

	reqStream, _, respStream, _, err := dump.Dump(id)
	defer func() {
		if reqStream != nil {
			reqStream.Close()
		}
		if respStream != nil {
			respStream.Close()
		}
	}()
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}

	if dumpType == "request" {
		respStream.Close()
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
		req, err := http.ReadRequest(bufio.NewReader(reqStream))
		if err != nil {
			response.Code = 1
			response.Message = err.Error()
			ctx.JSON(500, response)
			return
		}
		_, err = io.Copy(ctx.Writer, req.Body)
		if err != nil {
			response.Code = 1
			response.Message = err.Error()
			ctx.JSON(500, response)
			return
		}
	} else {
		reqStream.Close()
		resp, err := http.ReadResponse(bufio.NewReader(respStream), nil)
		if err != nil {
			response.Code = 1
			response.Message = err.Error()
			ctx.JSON(500, response)
			return
		}
		var r io.Reader
		if resp.Header.Get("Content-Encoding") == "gzip" {
			r, err = gzip.NewReader(resp.Body)
			if err != nil {
				log.Logger.Errorf("[Kasiopea-Controller] [%d] gzip init for response failed: %v", id, err)
				response.Code = 1
				response.Message = err.Error()
				ctx.JSON(500, response)
				return
			}
		} else if resp.Header.Get("Content-Encoding") == "deflate" {
			r, err = zlib.NewReader(resp.Body)
			if err != nil {
				log.Logger.Errorf("[Kasiopea-Controller] [%d] deflate init for response failed: %v", id, err)
				response.Code = 1
				response.Message = err.Error()
				ctx.JSON(500, response)
				return
			}
		} else {
			r = resp.Body
		}
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
		_, err = io.Copy(ctx.Writer, r)
		if err != nil {
			response.Code = 1
			response.Message = err.Error()
			ctx.JSON(500, response)
			return
		}
	}
}
