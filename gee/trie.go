package gee

import "strings"


type node struct{
	// specific string like go/:lang/123
	pattern string
	// certain parts// like go or :lang
	part string
	children []*node
	// wild card matching
	isWild bool
}

// match single node for insertion
func (n* node)matchChild(part string)*node{
	for _,child := range n.children{
		if child.part == part || child.isWild{
			return child
		}
	}
	return nil
}

// match mutiple node for search
func (n *node)matchChildren(part string)[]*node{
	nodes := make([]*node, 0)
	for _,child := range n.children{
		if(child.part == part || child.isWild){
			nodes = append(nodes,child)
		}
	}
	return nodes
}

// insert new route into trie tree
func(n *node)insert(pattern string, parts []string,height int){
	if len(parts) == height{
		n.pattern = pattern
		return 
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil{
		child = &node{part:part,isWild: part[0]==':'||part[0]=='*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern,parts,height+1)

}

// search route in the trie tree
func(n *node)search(parts[]string,height int)*node{
	if len(parts) == height || strings.HasPrefix(n.part,"*"){
		// check if we reach the last part or the last part is a wildcard
		if(n.pattern == ""){
			//match failed
			return nil
		}
		return n
	}
	// get current part and match all routes for current part
	part := parts[height]
	children := n.matchChildren(part)
	// for each part we recursively search for next part
	for _,child:=range children{
		result := child.search(parts,height+1)
		if(result != nil){
			// match found
			return result
		}

	}
	// match failed
	return nil

}





