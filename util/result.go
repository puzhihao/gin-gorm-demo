package util

import (
	"github.com/gin-gonic/gin"
	"github.com/puzhihao/gin-gorm-demo/types"
	"net/http"
	"time"
)

type Result struct {
	Time time.Time   `json:"time"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// 成功
func Success(c *gin.Context, data interface{}) {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Time = time.Now()
	res.Code = int(types.ApiCode.SUCCESS)
	res.Msg = types.ApiCode.GetMessage(types.ApiCode.SUCCESS)
	res.Data = data

	c.JSON(http.StatusOK, res)
}

// 出错
func MyError(c *gin.Context, msg string) {
	res := Result{}
	res.Time = time.Now()
	res.Code = int(types.ApiCode.FAILED)
	res.Msg = msg
	res.Data = gin.H{}

	c.JSON(http.StatusOK, res)
}
