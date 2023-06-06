package otz

import (
	"flag"
	"github.com/ShadowsGtt/otz/log"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

// Server 服务信息
type Server struct {
	engine *gin.Engine
}

func (s *Server) Register(method string, handler func(*gin.Context)) {
	s.engine.Any(method, handler)
}

func (s *Server) Start() error {
	cfg := GetGlobalConfig()
	if err := s.engine.Run(cfg.Ip + ":" + cfg.Port); err != nil {
		return err
	}
	return nil
}

func getServerConfigPath() string {
	// 如果是默认配置文件，则从命令行参数中获取
	if GlobalServerConfigFile == defaultConfigFile {
		flag.StringVar(&GlobalServerConfigFile, "conf", defaultConfigFile, "server config file")
		flag.Parse()
	}
	// 否则可以从全局变量中获取 - 允许用户设置,也便于单测
	return GlobalServerConfigFile
}

// NewServer 创建服务
func NewServer() *Server {
	s := &Server{}

	// 加载服务配置
	cfg, err := LoadConfig(getServerConfigPath())
	if err != nil {
		panic(err)
	}
	// 设置全局变量
	SetGlobalConfig(cfg)

	// 初始化日志
	log.SetDefaultLogger(cfg.Log)

	// 创建gin引擎
	gin.DefaultWriter = ioutil.Discard
	gin.DefaultErrorWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	s.engine = engine

	return s
}
