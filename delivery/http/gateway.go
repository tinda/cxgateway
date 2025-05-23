package http

import (
	"fmt"
	"github.com/codingXiang/configer"
	"github.com/codingXiang/go-logger"
	"github.com/codingXiang/gogo-i18n"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tinda/cxgateway/delivery"
	"github.com/tinda/cxgateway/middleware"
	http2 "github.com/tinda/cxgateway/module/auto_register/delivery/http"
	"github.com/tinda/cxgateway/module/auto_register/repository"
	"github.com/tinda/cxgateway/module/auto_register/service"
	"github.com/tinda/cxgateway/pkg/util"
	"net/http"
	"time"
)

type ApiGateway struct {
	engine      *gin.Engine
	Api         *gin.RouterGroup
	handler     util.RequestHandlerInterface
	configName  string
	uploadPath  string
	defaultData []byte
}

var (
	Gateway delivery.HttpHandler
)

func NewApiGatewayWithData(configName string, defaultData []byte) delivery.HttpHandler {
	var (
		gateway = &ApiGateway{}
	)
	//初始化 configer
	configer.Config = configer.NewConfiger()
	//設定多語系 Handler
	gogo_i18n.LangHandler = gogo_i18n.NewLanguageHandler()
	//設定預設資料
	configer.Config.AddCore(configName, configer.NewConfigerCore("yaml", ""))
	if data, err := configer.Config.GetCore(configName).ReadConfig(defaultData); err == nil {

		var (
			log  = data.Get("application.log")
			mode = data.GetString("application.mode")
		)

		//設定 log 等級與格式
		logger.Log = logger.NewLogger(logger.InterfaceToLogger(log))
		//伺服器模式
		logger.Log.Info("Server Mode =", mode)
		gin.SetMode(mode)

		gateway = &ApiGateway{
			engine:      gin.Default(),
			defaultData: defaultData,
			configName:  configName,
			uploadPath:  data.GetString("application.uploadPath"),
		}

		gateway.handler = util.NewRequestHandler()

		gateway.engine.
			Use(cors.Default()).
			Use(middleware.Logger(), gin.Recovery()).
			Use(middleware.RequestIDMiddleware(data.GetString("application.appId"))).
			Use(middleware.RequestVersion(data.GetString("application.version"))).
			Use(middleware.GoI18nMiddleware(data))
		gateway.Api = gateway.engine.Group(data.GetString("application.apiBaseRoute"))
	} else {
		panic(fmt.Sprintf("config %s is not set", configName))
	}

	return gateway
}

func NewApiGateway(configName string, core configer.CoreInterface) delivery.HttpHandler {
	var (
		//config  = settings.ConfigData.Data.Application
		gateway = &ApiGateway{}
	)
	if configName != "" {
		gateway.configName = configName
	} else {
		gateway.configName = "default"
	}
	//初始化 configer
	configer.Config = configer.NewConfiger()
	//設定多語系 Handler
	gogo_i18n.LangHandler = gogo_i18n.NewLanguageHandler()
	//設定預設資料
	if core == nil {
		if gateway.configName == "default" {
			gateway.defaultData = []byte(`application:
  timeout:
    read: 1000
    write: 1000
  port: 8080
  mode: "test"
  log:
    level: "debug"
    format: "json"
  appId: "app"
  appToken: ""
  apiBaseRoute: "/api"
i18n:
  defaultLanguage: "zh_Hant"
  file:
    path: "./i18n"
    type: "yaml"
`)
			configer.Config.AddCore(gateway.configName, configer.NewConfigerCore("yaml", ""))
		}
	} else {
		configer.Config.AddCore(gateway.configName, core)
	}
	if data, err := gateway.GetConfig().ReadConfig(gateway.defaultData); err == nil {
		var (
			//log          = data.Get("application.log")
			appId        = data.GetString("application.appId")
			appVersion   = data.GetString("application.version")
			apiBaseRoute = data.GetString("application.apiBaseRoute")
			mode         = data.GetString("application.mode")
			uploadPath   = data.GetString("application.uploadPath")
		)
		//設定 log 等級與格式
		//logger.Log = logger.NewLogger(logger.InterfaceToLogger(log))
		logger.Log = logger.NewLoggerWithConfiger(data)
		//伺服器模式
		logger.Log.Info("Server Mode =", mode)
		gin.SetMode(mode)
		gateway.handler = util.NewRequestHandler()
		gateway = &ApiGateway{
			engine:     gin.Default(),
			configName: configName,
			uploadPath: uploadPath,
		}

		//設定 cors
		if data.Get("cors") == nil {
			gateway.GetEngine().Use(cors.Default())
		} else {
			gateway.GetEngine().Use(middleware.GoCors(data))
		}

		//設定 gateway 的中間件
		gateway.engine.
			Use(middleware.Logger(), gin.Recovery()).
			Use(middleware.RequestIDMiddleware(appId)).
			Use(middleware.RequestVersion(appVersion)).
			Use(middleware.GoI18nMiddleware(data))
		gateway.Api = gateway.engine.Group(apiBaseRoute)
	} else {
		logger.Log.Error(fmt.Sprintf("config %s is not set, error = %s", gateway.configName, err))
	}

	return gateway
}

func (gateway *ApiGateway) GetEngine() *gin.Engine {
	return gateway.engine
}

func (gateway *ApiGateway) GetHandler() util.RequestHandlerInterface {
	return gateway.handler
}

func (gateway *ApiGateway) GetApiRoute() *gin.RouterGroup {
	return gateway.Api
}

func (this *ApiGateway) EnableAutoRegistration(configName string, configType string, configPath ...string) error {
	data := configer.NewConfigerCore(configType, configName, configPath...)
	data.SetAutomaticEnv("auto_registration")

	repo, err := repository.NewAutoRegisteredRepository(data)
	if err != nil {
		return err
	}
	repo.Initial()

	svc := service.NewAutoRegisteredService(repo)
	http2.NewAutoRegisteredHttpHandler(this, svc)

	return nil
}

func (this *ApiGateway) GetConfig() configer.CoreInterface {
	return configer.Config.GetCore(this.configName)
}

func (this *ApiGateway) GetUploadPath() string {
	return this.uploadPath
}

func (this *ApiGateway) Run() {
	if data, err := this.GetConfig().ReadConfig(this.defaultData); err == nil {
		var (
			port         = data.GetInt("application.port")          //伺服器的 port
			writeTimeout = data.GetInt("application.timeout.write") //伺服器的寫入超時時間
			readTimeout  = data.GetInt("application.timeout.read")  //伺服器讀取超時時間
		)
		logger.Log.Debug("Server port =", port)
		// 設定 http server
		server := &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        Gateway.GetEngine(),
			ReadTimeout:    time.Duration(readTimeout) * time.Second,
			WriteTimeout:   time.Duration(writeTimeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		logger.Log.Info("API Gateway Start Running")
		//啟動 http server
		server.ListenAndServe()
	}
}
