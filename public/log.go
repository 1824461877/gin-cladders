package public

import (
	"context"
	"gin-cladder/conf/elite/control"
	"github.com/gin-gonic/gin"
)

// 错误日志
func ContextWarning(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*control.TraceContext)
	if !ok{
		traceContext = control.NewTrace()
	}
	control.LogInfo.TagWarn(traceContext, dltag, m)
}

// 错误日志
func ContextError(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*control.TraceContext)
	if !ok{
		traceContext = control.NewTrace()
	}
	control.LogInfo.TagError(traceContext, dltag, m)
}

// 普通日志
func ContextNotice(c context.Context, dltag string, m map[string]interface{}) {
	v:=c.Value("trace")
	traceContext,ok := v.(*control.TraceContext)
	if !ok{
		traceContext = control.NewTrace()
	}
	control.LogInfo.TagInfo(traceContext, dltag, m)
}

// 错误日志
func ComLogWarning(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	control.LogInfo.TagError(traceContext, dltag, m)
}

//普通日志
func ComLogNotice(c *gin.Context, dltag string, m map[string]interface{}) {
	traceContext := GetGinTraceContext(c)
	control.LogInfo.TagInfo(traceContext, dltag, m)
}

// 从gin的Context中获取数据
func GetGinTraceContext(c *gin.Context) *control.TraceContext {
	// gin context 对其获取相关的数据
	if c == nil {
		return control.NewTrace()
	}
	traceContext, exists := c.Get("trace")
	if exists {
		if tc, ok := traceContext.(*control.TraceContext); ok {
			return tc
		}
	}
	return control.NewTrace()
}

// 从Context中获取数据
func GetTraceContext(c context.Context) *control.TraceContext {
	if c == nil {
		return control.NewTrace()
	}
	traceContext:=c.Value("trace")
	if tc, ok := traceContext.(*control.TraceContext); ok {
		return tc
	}
	return control.NewTrace()
}
