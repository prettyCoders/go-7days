package gee

import "strings"

/**
node 动态路由树节点
与普通的树不同，为了实现动态路由匹配，加上了isWild这个参数。
即当我们匹配 /p/go/doc/这个路由时，第一层节点，p精准匹配到了p，第二层节点，go模糊匹配到:lang，
那么将会把lang这个参数赋值为go，继续下一层匹配
*/
type node struct {
	routerPath string  //待匹配路由，例如 /p/:lang
	part       string  //路由中的一部分，例如 :lang
	childNodes []*node //子节点，例如 [doc, tutorial, intro]
	isWild     bool    //是否模糊匹配，part 含有 : 或 * 时为true
}

//matchChildNode 匹配子节点，返回匹配到的第一个子节点，用于插入
func (n *node) matchChildNode(part string) *node {
	for _, childNode := range n.childNodes {
		if childNode.part == part || childNode.isWild {
			return childNode
		}
	}
	return nil
}

//matchChildNodes 匹配子节点列表，返回所有匹配到到子节点，用于查找
func (n *node) matchChildNodes(part string) []*node {
	nodes := make([]*node, 0)
	for _, childNode := range n.childNodes {
		if childNode.part == part || childNode.isWild {
			nodes = append(nodes, childNode)
		}
	}
	return nodes
}

//insert 递归插入路由路径的每个部分
func (n *node) insert(routerPath string, parts []string, height int) {
	//递归结束条件
	if len(parts) == height {
		n.routerPath = routerPath
		return
	}

	part := parts[height]
	child := n.matchChildNode(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.childNodes = append(n.childNodes, child)
	}
	child.insert(routerPath, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.routerPath == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildNodes(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
