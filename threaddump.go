package actuator

import (
	"bytes"
	"errors"
	"net/http"
)

func getThreadDump() ([]byte, error) {
	var buffer bytes.Buffer
	profile := pprofLookupFunction(threadsKey)
	if profile == nil {
		return nil, errors.New(profileNotFoundError)
	}
	err := profile.WriteTo(&buffer, 1)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// handleThreadDump is the handler to get the thread dump
func handleThreadDump(writer http.ResponseWriter, requestRef *http.Request) {
	body, err := getThreadDump()
	HandleDump(body, err, writer, requestRef)
}
