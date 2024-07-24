package models

import (
	"errors"
	"log"
	"strings"
)

type NodeType string

const (
	ObjectNode  NodeType = "object"
	ArrayNode   NodeType = "array"
	StringNode  NodeType = "string"
	NumberNode  NodeType = "number"
	BooleanNode NodeType = "boolean"
	NullNode    NodeType = "null"
)

type JSONNode struct {
	Key      string               `json:"key"`
	Type     NodeType             `json:"type"`
	Size     int                  `json:"size,omitempty"`
	Children map[string]*JSONNode `json:"children,omitempty"`
	Value    interface{}          `json:"value,omitempty"`
}
type JSONTrie struct {
	Root *JSONNode
}

func NewJSONTrie() *JSONTrie {
	return &JSONTrie{
		Root: &JSONNode{
			Type:     ObjectNode,
			Children: make(map[string]*JSONNode),
		},
	}
}

func (t *JSONTrie) Insert(path string, value interface{}) {
	parts := strings.Split(path, ".")
	current := t.Root

	for i, part := range parts {
		if current.Children == nil {
			current.Children = make(map[string]*JSONNode)
		}

		if _, exists := current.Children[part]; !exists {
			current.Children[part] = &JSONNode{Key: part}
		}

		if i == len(parts)-1 {
			node := current.Children[part]
			node.Value = value
			node.Type = getNodeType(value)
			if node.Type == ArrayNode || node.Type == ObjectNode {
				node.Size = getSize(value)
			}
		}

		current = current.Children[part]
	}
}

func (t *JSONTrie) Get(path string) (*JSONNode, error) {
	if path == "" {
		log.Println("ROOT CASE")
		return t.Root, nil
	}
	parts := strings.Split(path, ".")
	current := t.Root
	log.Println("invoked GET method for JSONTrie", parts, current)

	for _, part := range parts {
		log.Println("starting with part", part)

		if current.Children == nil {
			log.Println("failed here")
			return nil, errors.New("path not found")
		}

		log.Println(current.Children)
		if child, exists := current.Children[part]; exists {
			current = child
		} else {
			log.Println("failed here 2")
			return nil, errors.New("path not found")
		}
	}

	return current, nil
}

func (t *JSONTrie) GetChildren(path string) ([]*JSONNode, error) {
	node, err := t.Get(path)
	if err != nil {
		return nil, err
	}

	if node.Children == nil {
		return nil, errors.New("node has no children")
	}

	children := make([]*JSONNode, 0, len(node.Children))
	for _, child := range node.Children {
		children = append(children, child)
	}

	return children, nil
}

func getNodeType(v interface{}) NodeType {
	switch v.(type) {
	case map[string]interface{}:
		return ObjectNode
	case []interface{}:
		return ArrayNode
	case string:
		return StringNode
	case float64, int, int64:
		return NumberNode
	case bool:
		return BooleanNode
	case nil:
		return NullNode
	default:
		return StringNode
	}
}

func getSize(v interface{}) int {
	switch val := v.(type) {
	case map[string]interface{}:
		return len(val)
	case []interface{}:
		return len(val)
	default:
		return 0
	}
}
