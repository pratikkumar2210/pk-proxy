package tree

import (
	"fmt"

	"github.com/pratikkumar2201/pk-proxy/models"
)

type TrieNode struct {
	IsEnd    bool
	Children map[string]*TrieNode
	Value    *models.Upstream
}

type RadixTree struct {
	Root *TrieNode
}

func NewRadixTree() *RadixTree {
	return &RadixTree{
		Root: &TrieNode{
			Children: make(map[string]*TrieNode),
		},
	}
}

// returns longest common prefix in both strings (index)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func lcp(s1, s2 string) int {
	n := min(len(s1), len(s2))
	i := 0
	for i < n && s1[i] == s2[i] {
		i++
	}
	return i
}

func (tree *RadixTree) Print(node *TrieNode, prefix string) {
	if node == nil {
		return
	}
	for edge, child := range node.Children {
		fmt.Printf("%s├── %s (end=%v) (val=%v)\n", prefix, edge, child.IsEnd, child.Value)
		tree.Print(child, prefix+"    ")
	}
}

func (tree *RadixTree) Insert(path string, value *models.Upstream) {
	node := tree.Root

	for {
		matched := false

		for edge, child := range node.Children {
			i := lcp(path, edge)

			if i == 0 {
				continue
			}

			matched = true

			// full edge match -> go deeper
			if i == len(edge) {
				node = child
				path = path[i:]
				goto NEXT
			}

			// split edge
			newNode := &TrieNode{
				IsEnd:    false,
				Children: make(map[string]*TrieNode),
				Value:    nil, // intermediate node has no value
			}

			// existing child becomes child of newNode
			newNode.Children[edge[i:]] = child

			// replace old edge
			node.Children[edge[:i]] = newNode
			delete(node.Children, edge)

			remaining := path[i:]

			if remaining == "" {
				// end of path
				newNode.IsEnd = true
				newNode.Value = value
			} else {
				// creatin new leaf node for remaining path
				newNode.Children[remaining] = &TrieNode{
					IsEnd:    true,
					Children: make(map[string]*TrieNode),
					Value:    value,
				}
			}
			return
		}

		// no match found -> insert directly
		if !matched {
			node.Children[path] = &TrieNode{
				IsEnd:    true,
				Children: make(map[string]*TrieNode),
				Value:    value,
			}
			return
		}

	NEXT:
		if path == "" {
			node.IsEnd = true
			node.Value = value
			return
		}
	}
}

func (tree *RadixTree) longestPrefix(route string) *TrieNode {
	node := tree.Root

	var bestNode *TrieNode

	path := ""

	for len(route) > 0 {
		var (
			bestEdge  string
			bestChild *TrieNode
			bestLCP   = 0
		)

		// find the edge with maximum LCP
		for edge, child := range node.Children {
			i := lcp(route, edge)
			if i > bestLCP {
				bestLCP = i
				bestEdge = edge
				bestChild = child
			}
		}

		// no match at all
		if bestLCP == 0 {
			break
		}

		// add matched part to path
		path += bestEdge[:bestLCP]

		// case 1: partial match inside edge → stop
		if bestLCP < len(bestEdge) {
			break
		}

		// case 2: full edge match → continue deeper
		route = route[bestLCP:]
		node = bestChild

		if node.IsEnd {
			bestNode = node
		}
	}

	return bestNode
}

func (tree *RadixTree) FindRoute(path string) *TrieNode {
	return tree.longestPrefix(path)
}
