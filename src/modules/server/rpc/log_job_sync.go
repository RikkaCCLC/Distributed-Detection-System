package rpc

import (
	"github.com/toolkits/pkg/logger"
	"open-devops/src/models"
)

// input = agent.hostName
func (*Server) LogJobSync(input string, output *[]*models.LogStrategy) error {
	ljs, _ := models.LogJobGets("id>0")

	*output = ljs
	logger.Infof("LogJobSync.call.receive ljs :%v %v %v", ljs, input, output)
	return nil

}
