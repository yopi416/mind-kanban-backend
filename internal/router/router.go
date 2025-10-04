package router

import (
	"net/http"

	"github.com/yopi416/mind-kanban-backend/internal/handler"
)

func NewRouter() *http.ServeMux {

	// register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.NewHealthzHandler().ServeHTTP)
	// mux.HandleFunc("/todos", handler.NewTODOHandler(svc).ServeHTTP)
	// mux.Handle("/do-panic", handler.NewDoPanicHandler())
	// mux.Handle("/do-panic-recover", middleware.Recovery(handler.NewDoPanicHandler()))
	// mux.Handle("/osdetect", middleware.OSDetection(handler.NewMockHandler()))
	// mux.Handle("/logging", middleware.OSDetection(middleware.Logging(handler.NewMockHandler())))
	// mux.Handle("/basicauth", middleware.BasicAuth(handler.NewMockHandler()))

	return mux
}
