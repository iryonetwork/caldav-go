package handlers

import (
	"net/http"

	"github.com/iryonetwork/caldav-go/data"
	"github.com/iryonetwork/caldav-go/global"
)

// HandlerInterface represents a CalDAV request handler. It has only one function `Handle`,
// which is used to handle the CalDAV request and returns the response.
type HandlerInterface interface {
	Handle() *Response
}

// Common data shared across the specific handlers. Defined here to
// easily make available, in a single place, all the basic data possibly needed by the handlers.
type handlerData struct {
	request     *http.Request
	requestBody string
	requestPath string
	headers     headers
	response    *Response
	storage     data.Storage
}

// NewHandler returns a new CalDAV request handler object based on the provided request.
// With the returned request handler, you can call `Handle()` to handle the request.
func NewHandler(request *http.Request) HandlerInterface {
	return NewHandlerWithStorage(request, global.Storage)
}

// NewHandlerWithStorage returns a new CalDAV request handler object based on the provided request.
// With the returned request handler, you can call `Handle()` to handle the request.
// Unlike NewHandler this function accepts an external storage and does not user the globally set storage.
func NewHandlerWithStorage(request *http.Request, storage data.Storage) HandlerInterface {
	hData := handlerData{
		request:     request,
		requestBody: readRequestBody(request),
		requestPath: request.URL.Path,
		headers:     headers{request.Header},
		response:    NewResponse(),
		storage:     storage,
	}

	switch request.Method {
	case "GET":
		return getHandler{handlerData: hData, onlyHeaders: false}
	case "HEAD":
		return getHandler{handlerData: hData, onlyHeaders: true}
	case "PUT":
		return putHandler{hData}
	case "DELETE":
		return deleteHandler{hData}
	case "PROPFIND":
		return propfindHandler{hData}
	case "OPTIONS":
		return optionsHandler{hData}
	case "REPORT":
		return reportHandler{hData}
	default:
		return notImplementedHandler{hData}
	}
}
