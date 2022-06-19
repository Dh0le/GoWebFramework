package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
}

// constuctor of gee.Engine
func New() *Engine{
	return &Engine{router: newRouter()}
}

// method that we use to add REST api
func (engine *Engine)addRoute(method string,pattern string,handler HandlerFunc){
	engine.router.addRoute(method,pattern,handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc){
	engine.addRoute("POST",pattern,handler)
}

func (engine *Engine) PUT(pattern string, handler HandlerFunc){
	engine.addRoute("PUT",pattern,handler)
}

func (engine *Engine) Run(addr string) (err error){
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	c := newContext(w,req)
	engine.router.handle(c)
}

