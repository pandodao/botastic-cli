package scan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pandodao/botastic-go"
)

func extractEntrieFile(file string) ([]*botastic.CreateIndexesItem, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content := strings.TrimSpace(string(bytes))

	if content == "" {
		return nil, nil
	}

	items := []*botastic.CreateIndexesItem{
		{
			ObjectID:   fmt.Sprintf("%s/%s-%d", filepath.Dir(file), filepath.Base(file), 0),
			Category:   "plain-text",
			Data:       content,
			Properties: "{}",
		},
	}

	return items, nil
}
