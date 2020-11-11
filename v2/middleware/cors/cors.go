package cors

import (
	"github.com/codingXiang/configer/v2"
	"github.com/codingXiang/cxgateway/v2/middleware"
	"github.com/codingXiang/cxgateway/v2/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strings"
)

type Cors struct {
	config *viper.Viper
}

func New(config *viper.Viper) middleware.Object {
	return &Cors{
		config: config,
	}

}

func (c *Cors) Handle() gin.HandlerFunc {
	if c.config.GetBool(server.Default_) {
		return cors.Default()
	} else {
		return cors.New(NewConfig(c.config))
	}
}

func NewConfig(data *viper.Viper) cors.Config {
	if data.GetStringMap(server.Cors) == nil {
		return cors.DefaultConfig()
	}
	allowAllOrigin := data.GetBool(configer.GetConfigPath(server.Cors, server.AllowAllOrigin))
	config := cors.DefaultConfig()
	if allowAllOrigin {
		config.AllowAllOrigins = allowAllOrigin
	} else {
		config.AllowOrigins = strings.Split(configer.GetConfigPath(server.Cors, server.AllowOrigins), ",")
	}
	config.AllowHeaders = strings.Split(configer.GetConfigPath(server.Cors, server.AllowHeaders), ",")
	config.AllowMethods = strings.Split(configer.GetConfigPath(server.Cors, server.AllowMethods), ",")
	return config
}