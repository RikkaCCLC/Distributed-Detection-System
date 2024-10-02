package web

import (
	"github.com/gin-gonic/gin"
	"open-devops/src/common"
	"open-devops/src/models"
)

func LogJobAdd(c *gin.Context) {

	var input models.LogStrategy
	if err := c.BindJSON(&input); err != nil {
		common.JSONR(c, 400, err)
		return
	}
	id, err := input.AddOne()
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}
	common.JSONR(c, 200, id)
}

func LogJobGets(c *gin.Context) {

	ljs, err := models.LogJobGets("id>0")
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, ljs)
}
