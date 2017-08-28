package echopprof

import (
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

type customEchoHandler struct {
	httpHandler http.Handler

	wrappedHandleFunc echo.HandlerFunc
	once              sync.Once
}

func (ceh *customEchoHandler) Handle(c echo.Context) error {
	ceh.once.Do(func() {
		ceh.wrappedHandleFunc = ceh.mustWrapHandleFunc(c)
	})
	return ceh.wrappedHandleFunc(c)
}

func (ceh *customEchoHandler) mustWrapHandleFunc(c echo.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		ceh.httpHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

func fromHTTPHandler(httpHandler http.Handler) *customEchoHandler {
	return &customEchoHandler{httpHandler: httpHandler}
}

func fromHandlerFunc(serveHTTP func(w http.ResponseWriter, r *http.Request)) *customEchoHandler {
	return &customEchoHandler{httpHandler: &customHTTPHandler{serveHTTP: serveHTTP}}
}
