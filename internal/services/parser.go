package services

import (
	"encoding/json"
	"github.com/shammianand/fast-json-viewer/internal/models"
	"io"
)

type Parser struct {
	MaxFileSize int64
}

func NewParser(maxFileSize int64) *Parser {
	return &Parser{MaxFileSize: maxFileSize}
}

func (p *Parser) ParseJSON(r io.Reader) (*models.JSONTrie, error) {
	trie := models.NewJSONTrie()
	decoder := json.NewDecoder(r)

	var data interface{}
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	p.parseValue("", data, trie)

	return trie, nil
}

func (p *Parser) parseValue(path string, v interface{}, trie *models.JSONTrie) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, v := range val {
			newPath := path
			if newPath != "" {
				newPath += "."
			}
			newPath += k
			p.parseValue(newPath, v, trie)
		}
	case []interface{}:
		for i, v := range val {
			newPath := path + "[" + string(i) + "]"
			p.parseValue(newPath, v, trie)
		}
	default:
		trie.Insert(path, v)
	}
}
