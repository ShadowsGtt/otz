package log_test

import (
	"context"
	"errors"
	"github.com/ShadowsGtt/otz/log"
	"github.com/ShadowsGtt/otz/otzctx"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestSetDefaultLogger(t *testing.T) {
	// 默认控制台输出测试
	log.Infof("this will output to console")

	// 重置
	log.SetDefault([]log.Config{
		{
			OutputType: log.OutputTypeFile,
			FileName:   "test.log",
			Level:      "debug",
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 1,
			Compress:   false,
			FormatType: log.FormatTypeText,
		},
		{
			OutputType: log.OutputTypeFile,
			FileName:   "test1.log",
			Level:      "info",
			MaxSize:    1,
			MaxAge:     1,
			MaxBackups: 1,
			Compress:   false,
			FormatType: log.FormatTypeText,
		},
	}...)
	log.Infof("this will output to console")
	log.Debugf("output to file, this log level is debug, output to test.log, not output to test1.log")

	// 重置到控制台
	log.SetDefault(log.Config{OutputType: log.OutputTypeConsole})
}

func TestWithCtx(t *testing.T) {
	otzCtx := otzctx.GetOrNewOTZContext(context.Background())
	defer otzctx.PutOTZCtx(otzCtx)
	otzCtx.SetLogger(log.NewZapLog(log.Config{OutputType: log.OutputTypeConsole}))
	ctx := otzCtx.Context()
	log.WithCtx(ctx, "age", "18")
	log.WithCtx(ctx, "name", "tom")

	log.InfoCtxf(ctx, "===========")
}

func TestWith(t *testing.T) {
	logger := log.NewZapLog(log.Config{OutputType: log.OutputTypeConsole})
	logger = logger.With("name", "tom")
	logger = logger.With("age", "18")

	logger.Infof("1111")
	logger.Infof("2222")
}

type YamlParser struct {
	node *yaml.Node
}

// Parse 解析yaml
func (y *YamlParser) Parse(dst interface{}) error {
	if y.node == nil {
		return errors.New("node is nil")

	}
	return y.node.Decode(dst)
}

func TestSetDefaultWithLoader(t *testing.T) {
	data := `
log:
  - output_type: file # 日志输出 console/file
    file_name: ./loader_test.log   # 文件名
    level: debug # 日志级别 debug info warn error fatal
    max_size: 10   # 文件大小限制 MB
    max_age: 7    # 保留天数
    max_backups: 10 # 文件数
    compress: false    # 是否压缩
    format_type: text # 格式换类型 text/file
    
  - output_type: console  # 输出到控制台
    level: debug
    format_type: text
`

	type logConfig struct {
		Node yaml.Node `yaml:"log"`
	}

	lc := &logConfig{}
	err := yaml.Unmarshal([]byte(data), lc)
	if err != nil {
		t.Fatal(err)
	}

	yamlParse := &YamlParser{node: &lc.Node}
	err = log.SetDefaultWithLoader(yamlParse)
	if err != nil {
		t.Fatal(err)
	}
	otzCtx := otzctx.NewOtzContext(context.Background())
	defer otzctx.PutOTZCtx(otzCtx)
	ctx := otzCtx.Context()

	log.WithCtx(ctx, "name", "tom")
	log.WithCtx(ctx, "age", "18")
	log.InfoCtxf(ctx, "this test for parser")
}
