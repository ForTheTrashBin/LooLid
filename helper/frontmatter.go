package helper

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

// SerializeFrontMatter writes frontmatter content as a TOML-compatible block.
func SerializeFrontMatter(frontMatter map[string]any) ([]byte, error) {
	if len(frontMatter) == 0 {
		return []byte("+++\n+++\n"), nil
	}

	encoded, err := toml.Marshal(frontMatter)
	if err != nil {
		return nil, err
	}

	return []byte("+++\n" + string(encoded) + "+++\n"), nil
}

// DeserializeFrontMatter parses frontmatter content from a serialized block or raw TOML content.
func DeserializeFrontMatter(data []byte) (map[string]any, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return map[string]any{}, nil
	}

	content := trimmed
	lines := strings.Split(trimmed, "\n")
	if len(lines) >= 2 {
		firstLine := strings.TrimSpace(lines[0])
		if firstLine == "+++" || firstLine == "---" {
			closingIndex := -1
			for idx := 1; idx < len(lines); idx++ {
				if strings.TrimSpace(lines[idx]) == firstLine {
					closingIndex = idx
					break
				}
			}

			if closingIndex == -1 {
				return nil, fmt.Errorf("unterminated frontmatter block")
			}

			content = strings.Join(lines[1:closingIndex], "\n")
		}
	}

	frontMatter := map[string]any{}
	if _, err := toml.Decode(content, &frontMatter); err != nil {
		return nil, err
	}

	for key, value := range frontMatter {
		frontMatter[key] = normalizeFrontMatterValue(value)
	}

	return frontMatter, nil
}

func normalizeFrontMatterValue(value any) any {
	switch typedValue := value.(type) {
	case int64:
		return int(typedValue)
	case int32:
		return int(typedValue)
	case int16:
		return int(typedValue)
	case int8:
		return int(typedValue)
	case uint64:
		return int(typedValue)
	case uint32:
		return int(typedValue)
	case uint16:
		return int(typedValue)
	case uint8:
		return int(typedValue)
	case []any:
		normalized := make([]any, len(typedValue))
		for idx, item := range typedValue {
			normalized[idx] = normalizeFrontMatterValue(item)
		}
		return normalized
	case map[string]any:
		normalized := make(map[string]any, len(typedValue))
		for key, item := range typedValue {
			normalized[key] = normalizeFrontMatterValue(item)
		}
		return normalized
	default:
		return value
	}
}

// SerializeFrontMatterToString is a convenience wrapper for string-based callers.
func SerializeFrontMatterToString(frontMatter map[string]any) (string, error) {
	data, err := SerializeFrontMatter(frontMatter)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimRight(data, "\n")), nil
}
