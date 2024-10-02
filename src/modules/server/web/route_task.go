package web

import (
	"github.com/gin-gonic/gin"
	"open-devops/src/common"
	"open-devops/src/models"
	"strconv"
)

func TaskAdd(c *gin.Context) {

	var input models.TaskMeta
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

func TaskGets(c *gin.Context) {

	ljs, err := models.TaskMetaGets("id>0")
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, 200, ljs)
}

func TaskKill(c *gin.Context) {

	taskId, err := strconv.Atoi(c.DefaultQuery("task_id", "0"))
	if err != nil || taskId == 0 {
		common.JSONR(c, 400, err)
		return
	}
	err = models.SetTaskKill(int64(taskId))
	if err != nil {
		common.JSONR(c, 500, err)
		return
	}

	common.JSONR(c, "success")
}
