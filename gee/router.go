package gee

import (
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	roots map[string]*node
}

func newRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string)[]string{
	// only one * is allow in this wildcard matching and parsing
	vs := strings.Split(pattern,"/")
	parts := make([]string,0)
	for _,item := range vs{
		if item != ""{
			parts = append(parts, item)
			if(item[0]=='*'){
				break;
			}
		}
	}
	return parts
}



func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//first we split pattern into parts
	parts := parsePattern(pattern)
	// get key for current handler function
	key := method + "-" + pattern
	// check if current method is supported
	_,ok := r.roots[method]
	if(!ok){
		// add new method if current method is not supported
		r.roots[method] = &node{}
	}
	// insert new routes for this method
	r.roots[method].insert(pattern,parts,0)
	// insert handler into hashtable
	r.handlers[key] = handler
}

func (r *router)getRoute(method string, path string)(*node,map[string]string){
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root,ok := r.roots[method]

	if(!ok){
		// method not supported
		return nil,nil;
	}
	n := root.search(searchParts,0)

	if n != nil{
		// we found a match path
		// now we need to parse the parameter
		parts := parsePattern(n.pattern)
		for index, part := range parts{
			if part[0] == ':'{
				// found a paramter match
				params[part[1:]] =  searchParts[index]

			}else if part[0] == '*' && len(part) > 1{
				// found a wild card matching
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n,params
	}
	return nil,nil
}


func(r *router)handle(c *Context){
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		// since we could have alot of middleware before we just append it to the end
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c*Context){
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// execute in sequence 
	c.Next()
}


