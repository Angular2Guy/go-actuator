package actuator

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmptyConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Failf(t, "incorrect config", "%+v", r)
		}
	}()
	c := &Config{}
	c.validate()
}

func TestValidateConfigWithIncorrectEndpoint(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				assert.Equal(t, fmt.Sprintf("invalid endpoint %d provided", 20), e.Error())
			}
		}
	}()
	c := &Config{
		Endpoints: []int{20},
	}
	c.validate()
}

func TestSetDefaultsInConfig(t *testing.T) {
	c := &Config{}
	c.setDefaults()
	assert.Equal(t, AllEndpoints, c.Endpoints)
}

func TestEnv(t *testing.T) {
	w := setupMuxAndGetResponse(t, Env, envEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data)
}

func TestEnvNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, envEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestEnvInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{Env}}, http.MethodHead, envEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestEnvWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, envEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]string
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data)
}

func TestInfo(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, infoEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]interface{}
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data)
}

func TestInfoNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Env, infoEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestInfoInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{Info}}, http.MethodHead, infoEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestInfoWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, infoEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data map[string]interface{}
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data)
}

func TestMetrics(t *testing.T) {
	w := setupMuxAndGetResponse(t, Metrics, metricsEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data MetricsResponse
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data.MemStats)
}

func TestMetricsNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, metricsEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestMetricsInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{Metrics}}, http.MethodHead, metricsEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestMetricsWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, metricsEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	var data MetricsResponse
	getTypedJSONBody(t, w.Body, &data)
	assert.NotEmpty(t, data.MemStats)
}

func TestMetricsEncodeJSONError(t *testing.T) {
	mockEncodeJSONWithError()
	defer unMockEncodeJSON()

	w := setupMuxAndGetResponse(t, Metrics, metricsEndpoint)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error", w.Body.String())
	assert.Equal(t, textStringContentType, w.Header().Get(contentTypeHeader))
}

func TestPing(t *testing.T) {
	w := setupMuxAndGetResponse(t, Ping, pingEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Empty(t, w.Body)
}

func TestPingNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, pingEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestPingInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{Ping}}, http.MethodHead, pingEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestPingWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, pingEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Empty(t, w.Body)
}

/*
func TestShutdown(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing here, just to handle shutdown gracefully
		}
	}()
	w := setupMuxAndGetResponse(t, Shutdown, shutdownEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Empty(t, w.Body)
}

func TestShutdownNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, shutdownEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestShutdownInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{Shutdown}}, http.MethodHead, shutdownEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestShutdownWithoutConfig(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// do nothing here, just to handle shutdown gracefully
		}
		}()
		w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, shutdownEndpoint)
		assert.Equal(t, http.StatusOK, w.Code)

		assert.Empty(t, w.Body)
	}
*/

func TestThreadDump(t *testing.T) {
	w := setupMuxAndGetResponse(t, ThreadDump, threadDumpEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, w.Body)
}

func TestThreadDumpNotConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, threadDumpEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestThreadDumpInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{ThreadDump}}, http.MethodHead, threadDumpEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestThreadDumpWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, threadDumpEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, w.Body)
}

func TestThreadDumpWithError(t *testing.T) {
	mockPprofLookupWithError()
	defer unMockPprofLookup()

	w := setupMuxAndGetResponse(t, ThreadDump, threadDumpEndpoint)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	assert.Equal(t, profileNotFoundError, w.Body.String())
}

func TestGoRoutineDump(t *testing.T) {
	w := setupMuxAndGetResponse(t, GoRoutineDump, goRoutineDumpEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, w.Body)
}

func TestGoRoutineDumpConfigured(t *testing.T) {
	w := setupMuxAndGetResponse(t, Info, goRoutineDumpEndpoint)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, notFoundError, w.Body.String())
}

func TestGoRoutineDumpInvalidMethod(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, &Config{Endpoints: []int{ThreadDump}}, http.MethodHead, goRoutineDumpEndpoint)
	assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	assert.Equal(t, methodNotAllowedError, w.Body.String())
}

func TestGoRoutineDumpWithoutConfig(t *testing.T) {
	w := setupMuxWithConfigAndGetResponseForMethod(t, nil, http.MethodGet, goRoutineDumpEndpoint)
	assert.Equal(t, http.StatusOK, w.Code)

	assert.NotEmpty(t, w.Body)
}

func TestGoRoutineDumpWithError(t *testing.T) {
	mockPprofLookupWithError()
	defer unMockPprofLookup()

	w := setupMuxAndGetResponse(t, ThreadDump, threadDumpEndpoint)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	assert.Equal(t, profileNotFoundError, w.Body.String())
}

func TestGetLastStringAfterDelim(t *testing.T) {
	assert.Equal(t, "", getLastStringAfterDelimiter("", slash))
	assert.Equal(t, "a", getLastStringAfterDelimiter("a", slash))
	assert.Equal(t, "", getLastStringAfterDelimiter("", ""))
	assert.Equal(t, "c", getLastStringAfterDelimiter("a/b/c", slash))
	assert.Equal(t, "", getLastStringAfterDelimiter("a/", slash))
}
