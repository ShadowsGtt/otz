package otz

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/ShadowsGtt/otz/log"
	"github.com/ShadowsGtt/otz/otzctx"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"time"
)

// Server 服务信息
type Server struct {
	engine *gin.Engine
}

// Register 注册服务
func (s *Server) Register(method string, handler func(ctx context.Context)) {
	h := func(ginCtx *gin.Context) {
		otzCtx := otzctx.GetOrNewOTZContext(context.Background())
		otzCtx.SetGinCtx(ginCtx)
		begin := time.Now()
		defer func() {
			if err := recover(); err != nil {
				log.ErrorCtxf(otzCtx.Context(), "%s", string(debug.Stack()))
				ginCtx.Data(http.StatusInternalServerError, "", nil)
			}
			log.InfoCtxf(otzCtx.Context(), "URI: %s, cost: %dms",
				ginCtx.Request.URL.Path, time.Since(begin).Milliseconds(),
			)
		}()
		defer otzctx.PutOTZCtx(otzCtx)
		handler(otzCtx.Context())
	}
	s.engine.Any(method, h)
}

// Start 启动服务
func (s *Server) Start() error {
	cfg := GetGlobalConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Ip, cfg.Server.Port)
	log.Infof("server start, listen addr: %s", addr)
	if err := s.engine.Run(addr); err != nil {
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

type logParser struct {
	node *yaml.Node
}

// Parse 解析日志配置
func (p *logParser) Parse(dst interface{}) error {
	if p.node == nil {
		return errors.New("node is nil")
	}
	return p.node.Decode(dst)
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
	parser := &logParser{node: &cfg.Log}
	err = log.SetDefaultWithLoader(parser)
	if err != nil {
		panic(errors.New("parse log config failed, err: " + err.Error()))
	}

	// 创建gin引擎
	gin.DefaultWriter = ioutil.Discard
	gin.DefaultErrorWriter = ioutil.Discard
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	s.engine = engine

	return s
}
