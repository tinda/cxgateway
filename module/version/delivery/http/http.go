package http

import (
	"github.com/gin-gonic/gin"
	cx "github.com/tinda/cxgateway/delivery"
	"github.com/tinda/cxgateway/module/version"
	"github.com/tinda/cxgateway/module/version/delivery"
	"github.com/tinda/cxgateway/pkg/e"
	"github.com/tinda/cxgateway/pkg/i18n"
	"github.com/tinda/cxgateway/pkg/util"
)

const (
	MODULE = "version"
)

type VersionHttpHandler struct {
	i18nMsg i18n.I18nMessageHandlerInterface
	gateway cx.HttpHandler
	svc     version.Service
}

func NewVersionHttpHandler(gateway cx.HttpHandler, svc version.Service) delivery.HttpHandler {
	var handler = &VersionHttpHandler{
		i18nMsg: i18n.NewI18nMessageHandler(MODULE),
		gateway: gateway,
		svc:     svc,
	}
	/*
		v1 版本的 Ticket API
	*/
	v1 := gateway.GetEngine().Group("")
	v1.GET("", e.Wrapper(handler.GetServerVersion))
	//v1.GET("/check", e.Wrapper(handler.CheckVersion))
	//v1.POST("/upgrade", e.Wrapper(handler.Upgrade))
	return handler
}

func (this VersionHttpHandler) GetServerVersion(c *gin.Context) error {
	this.i18nMsg.SetCore(util.GetI18nData(c))
	if version, err := this.svc.GetServerVersion(); err != nil {
		return this.i18nMsg.GetError(err)
	} else {
		c.JSON(this.i18nMsg.GetSuccess(version))
	}
	return nil
}

func (this VersionHttpHandler) CheckVersion(c *gin.Context) error {
	this.i18nMsg.SetCore(util.GetI18nData(c))
	if err := this.svc.CheckVersion(); err != nil {
		return this.i18nMsg.GetError(err)
	} else {
		c.JSON(this.i18nMsg.GetSuccess(nil))
	}
	return nil
}

func (this VersionHttpHandler) Upgrade(c *gin.Context) error {
	this.i18nMsg.SetCore(util.GetI18nData(c))
	if err := this.svc.Upgrade(); err != nil {
		return this.i18nMsg.UpdateError(err)
	} else {
		c.JSON(this.i18nMsg.UpdateSuccess(nil))
	}
	return nil
}
