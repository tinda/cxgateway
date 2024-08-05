package delivery

import (
	"github.com/codingXiang/configer"
	"github.com/gin-gonic/gin"
	"github.com/tinda/cxgateway/pkg/util"
)

type HttpHandler interface {
	GetEngine() *gin.Engine
	GetApiRoute() *gin.RouterGroup
	GetHandler() util.RequestHandlerInterface
	GetConfig() configer.CoreInterface
	GetUploadPath() string
	EnableAutoRegistration(configName string, configType string, configPath ...string) error
	Run()
}
