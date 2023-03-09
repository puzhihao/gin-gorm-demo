package http

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/puzhihao/gin-gorm-demo/app/host"
	"github.com/puzhihao/gin-gorm-demo/app/host/impl"
	"github.com/puzhihao/gin-gorm-demo/util"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var (
	// BodyMaxContenxLength body最大大小 默认64M
	BodyMaxContenxLength int64 = 1 << 26
)

func GetHost(c *gin.Context) {
	var (
		PageSize   = 1
		PageNumber = 20
	)
	psStr := c.Query("PageSize")
	if psStr != "" {
		PageSize, _ = strconv.Atoi(psStr)
	}
	pnStr := c.Query("PageNumber")
	if pnStr != "" {
		PageNumber, _ = strconv.Atoi(pnStr)
	}
	keyWords := c.Query("KeyWords")
	req := &host.QueryHostRequest{
		PageSize:   PageSize,
		PageNumber: PageNumber,
		Keywords:   keyWords,
	}

	set, err := impl.Service.QueryHost(c.Request.Context(), req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
	} else {
		util.Success(c, set)
	}

}

func CreateHost(c *gin.Context) {

	req := host.NewDefaultHost()

	err := GetDataFromRequest(c.Request, req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))

	}
	ins, err := impl.Service.CreateHost(c.Request.Context(), req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
	} else {
		util.Success(c, ins)
	}
}

func DescribeHost(c *gin.Context) {
	req := &host.DescribeHostRequest{
		Id: c.Params.ByName("id"),
	}

	ins, err := impl.Service.DescribeHost(c.Request.Context(), req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
	} else {
		util.Success(c, ins)
	}

}

func UpdateHost(c *gin.Context) {
	req := host.NewPutUpdateHostRequest()

	err := GetDataFromRequest(c.Request, req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
		return
	}
	req.Id = c.Params.ByName("id")
	ins, err := impl.Service.UpdateHost(c.Request.Context(), req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
	} else {
		util.Success(c, ins)
	}
}

func PatchHost(c *gin.Context) {
	req := host.NewPatchUpdateHostRequest()

	err := GetDataFromRequest(c.Request, req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
		return
	}
	req.Id = c.Params.ByName("id")
	ins, err := impl.Service.UpdateHost(c.Request.Context(), req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
	} else {
		util.Success(c, ins)
	}
}

func DeleteHost(c *gin.Context) {
	id := c.Params.ByName("id")
	req := host.DeleteHostRequest{
		Id: id,
	}
	ins, err := impl.Service.DeleteHost(c.Request.Context(), &req)
	if err != nil {
		util.MyError(c, fmt.Sprintf("%s", err))
		log.Println(err)
	} else {
		util.Success(c, ins)
	}

}

func GetDataFromRequest(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}
