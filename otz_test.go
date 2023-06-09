package otz

import (
	"context"
	"github.com/ShadowsGtt/otz/log"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	GlobalServerConfigFile = "./test.yaml"
	s := NewServer()
	logCtx := log.NewLogCtx(context.Background())
	ctx := logCtx.GetCtx()
	log.WithCtx(ctx, "method", "TestNewServer")

	log.Debugf("server start")
	log.InfoCtxf(ctx, "server start...")

	s.Register("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		log.Infof("server request, time: %s", time.Now().Format("2006-01-02 15:04:05"))
	})
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
