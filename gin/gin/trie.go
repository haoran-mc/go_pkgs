package gin

import "strings"

type node struct {
	pattern  string  // route to be matched, such as /p/:lang
	part     string  // the current part of the route, such as :/lang
	children []*node // child nodes, such as [doc, tutorial, intro]
	isWild   bool    // Whether an exact match, The value is true when the route contains : or *
}

// matchChild finds the matched child for insert
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// matchChildren finds all the matched children for search
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// Recurse to the deepest point
	if height == len(parts) {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// Now we want to enter the page: /assets/css/style.css
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		// We use the field n.pattern to indicate whether the current node is the end of a full path
		// A match is successful only if the path of the current node is a complete path
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
