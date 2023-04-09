package actuator

import (
	"net/http"
	"runtime/pprof"
)

// variables for mocking in case of unit testing
var (
	encodeJSONFunction  = encodeJSON
	pprofLookupFunction = pprof.Lookup
)

func HandleDump(body []byte, err error, writer http.ResponseWriter, _ *http.Request) {
	if err != nil {
		// some error occurred
		// send the error in the response
		sendStringResponse(writer, http.StatusInternalServerError, err.Error())
		return
	}
	// now once we have the correct response
	writer.Header().Add(contentTypeHeader, textStringContentType)
	_, _ = writer.Write(body)
}
