package otzctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"sync"
)

const (
	ContextKey = "OTZ_CONTEXT_KEY"
)

var ctxPool = sync.Pool{
	New: func() interface{} {
		return &otzContext{}
	},
}

// Context 上下文
type Context interface {
	SetLogger(interface{})
	GetLogger() interface{}
	SetGinCtx(*gin.Context)
	GetGinCtx() *gin.Context
	Context() context.Context
}

type otzContext struct {
	logger  interface{}
	context context.Context
	ginCtx  *gin.Context
}

// SetLogger 设置logger
func (ctx *otzContext) SetLogger(logger interface{}) {
	ctx.logger = logger
}

// GetLogger 获取logger
func (ctx *otzContext) GetLogger() interface{} {
	return ctx.logger
}

// SetGinCtx 设置gin ctx
func (ctx *otzContext) SetGinCtx(ginCtx *gin.Context) {
	ctx.ginCtx = ginCtx
}

// GetGinCtx 获取gin ctx
func (ctx *otzContext) GetGinCtx() *gin.Context {
	return ctx.ginCtx
}

// Context 获取context
func (ctx *otzContext) Context() context.Context {
	return ctx.context
}

// OTZContext 转化成otz ctx
func OTZContext(ctx context.Context) Context {
	otzCtx, ok := ctx.Value(ContextKey).(*otzContext)
	if ok && otzCtx != nil {
		return otzCtx
	}
	return &otzContext{context: ctx}
}

// NewOtzContext 新建otz ctx
func NewOtzContext(ctx context.Context) Context {
	m := ctxPool.Get().(*otzContext)
	ctx = context.WithValue(ctx, ContextKey, m)
	m.context = ctx
	return m
}

// GetOrNewOTZContext 获取otz ctx,不存在则创建
func GetOrNewOTZContext(ctx context.Context) Context {
	// 从ctx中获取
	otzCtx, ok := ctx.Value(ContextKey).(*otzContext)
	if ok && otzCtx != nil {
		return otzCtx
	}
	return NewOtzContext(ctx)
}

// PutOTZCtx 回收otz ctx
func PutOTZCtx(ctx Context) {
	v, ok := ctx.(*otzContext)
	if !ok {
		return
	}
	v.logger = nil
	v.ginCtx = nil
	v.logger = nil
	ctxPool.Put(ctx)
}
