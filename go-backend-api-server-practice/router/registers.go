package router

import "net/http"

type RouterRegister struct {
	mux *CustomServerMux
}

func NewRouteRegister(mux *CustomServerMux) *RouterRegister {
	return &RouterRegister{mux: mux}
}

func (rr *RouterRegister) Register(method, pattern string, handler http.HandlerFunc, errMessage map[int]string) {
	rr.mux.HandleFunc(method+" "+pattern, handler)
	for statusCode, msg := range errMessage {
		rr.mux.RegisterErrorMessage(pattern, statusCode, msg)
	}
}
