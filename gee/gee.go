package gee

import (
	"log"
	"net/http"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
	//we have a parent routergroup for all sub routergroups
	*RouterGroup
	groups []*RouterGroup
}

type RouterGroup struct{
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	engine *Engine
}



// constuctor of gee.Engine
func New() *Engine{
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

//Group is defined to create a new RouterGroup
func (group *RouterGroup)Group(prefix string)*RouterGroup{
	// all Router group share a same instance of engine(singleton)
	engine := group.engine
	newGroup := &RouterGroup{
		prefix:group.prefix + prefix,
		parent: group,
		engine:engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}



// method that we use to add REST api
func (group *RouterGroup)addRoute(method string,comp string,handler HandlerFunc){
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s",method,pattern)
	group.engine.router.addRoute(method,pattern,handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc){
	group.addRoute("POST",pattern,handler)
}

func (group *RouterGroup) PUT(pattern string, handler HandlerFunc){
	group.addRoute("PUT",pattern,handler)
}

func (engine *Engine) Run(addr string) (err error){
	return http.ListenAndServe(addr,engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request){
	c := newContext(w,req)
	engine.router.handle(c)
}

