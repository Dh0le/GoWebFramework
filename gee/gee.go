package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

// HandlerFunc defines the request handler used by gee
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	router *router
	//the routergroup pointer will give engine all function of a routergroup
	*RouterGroup
	groups []*RouterGroup

	// add support for html template
	htmlTemplates *template.Template
	funcMap template.FuncMap
}

type RouterGroup struct{
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	// router group can access function of engine as singleton
	engine *Engine
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

//Use is define to add new middleware into the group

func(group *RouterGroup)Use(middlewares ...HandlerFunc){
	group.middlewares = append(group.middlewares, middlewares...)
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
	var middlewares []HandlerFunc
	for _,group:=range engine.groups{
		if strings.HasPrefix(req.URL.Path,group.prefix){
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w,req)
	c.handlers  = middlewares
	engine.router.handle(c)
}

// create static handler
func (group *RouterGroup)createStaticHandler(relativePath string ,fs http.FileSystem)HandlerFunc{
	absolutePath := path.Join(group.prefix,relativePath)
	fileServer := http.StripPrefix(absolutePath,http.FileServer(fs))
	return func(c *Context){
		file := c.Param("filepath")
		// check if we have that file or if we have the permission to access it
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer,c.Req)

	}
}

// serve static files
func (group *RouterGroup)Static(relativePath string, root string){
	handler := group.createStaticHandler(relativePath,http.Dir(root))
	urlPattern := path.Join(relativePath,"/*filepath")
	// register a get handler
	group.GET(urlPattern,handler)
}
