package otz

import (
	"github.com/ShadowsGtt/otz/log"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	GlobalServerConfigFile = "./test.yaml"
	s := NewServer()
	log.Debug("server start")

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
