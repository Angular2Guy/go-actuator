package core_test

import (
	"github.com/gin-gonic/gin"
	"github.com/sinhashubham95/go-actuator/commons"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sinhashubham95/go-actuator/core"
)

func TestGetHTTPTrace(t *testing.T) {
	assert.Empty(t, core.GetHTTPTrace())
}

func TestWithGIN(t *testing.T) {
	router := setupGINRouter()
	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)
	router.ServeHTTP(w, request)
	assert.Equal(t, http.StatusOK, w.Code)

	traces := core.GetHTTPTrace()
	assert.NotEmpty(t, traces)
	trace := traces[0]
	assert.Equal(t, http.MethodGet, trace.Request.Method)
	assert.Equal(t, "/ping", trace.Request.URL)
	assert.Empty(t, trace.Request.Headers)
	assert.Equal(t, http.StatusOK, trace.Response.Status)
	assert.Empty(t, trace.Response.Headers)
}

func TestWithNetHTTP(t *testing.T) {
	m := &http.ServeMux{}
	m.HandleFunc("/ping", getNetHTTPHandler())
	w := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/ping", nil)
	assert.NoError(t, err)
	m.ServeHTTP(w, request)
	assert.Equal(t, http.StatusOK, w.Code)

	traces := core.GetHTTPTrace()
	assert.NotEmpty(t, traces)
	trace := traces[0]
	assert.Equal(t, http.MethodGet, trace.Request.Method)
	assert.Equal(t, "/ping", trace.Request.URL)
	assert.Empty(t, trace.Request.Headers)
	assert.Equal(t, http.StatusOK, trace.Response.Status)
	assert.Empty(t, trace.Response.Headers)
}

func TestForMoreThanThresholdRequests(t *testing.T) {
	router := setupGINRouter()
	for i := 0; i <= (commons.HTTPTraceResultsSize + 5); i += 1 {
		w := httptest.NewRecorder()
		request, err := http.NewRequest(http.MethodGet, "/ping", nil)
		assert.NoError(t, err)
		router.ServeHTTP(w, request)
		assert.Equal(t, http.StatusOK, w.Code)
	}

	traces := core.GetHTTPTrace()
	assert.Equal(t, commons.HTTPTraceResultsSize, len(traces))
	trace := traces[0]
	assert.Equal(t, http.MethodGet, trace.Request.Method)
	assert.Equal(t, "/ping", trace.Request.URL)
	assert.Empty(t, trace.Request.Headers)
	assert.Equal(t, http.StatusOK, trace.Response.Status)
	assert.Empty(t, trace.Response.Headers)
}

func getFastHTTPHandler() fasthttp.RequestHandler {
	return core.WrapFastHTTPHandler(func(ctx *fasthttp.RequestCtx) {
		ctx.SetStatusCode(http.StatusOK)
	})
}

func setupGINRouter() *gin.Engine {
	router := gin.Default()
	router.Use(core.GINTracer())
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.Status(http.StatusOK)
	})
	return router
}

func getNetHTTPHandler() http.HandlerFunc {
	return core.WrapNetHTTPHandler(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
}