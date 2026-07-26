package nutsandbolts

import (
	"bytes"
	"strings"

	"github.com/BurntSushi/toml"
)

//-----------------------------------------------------------------------------
// SerializeFrontMatter writes frontmatter content as a TOML-compatible block.
//-----------------------------------------------------------------------------

func SerializeFrontMatter(frontMatter map[string]any) ([]byte, error) {

	retVal, err := []byte{}, error(nil)

	if len(frontMatter) > 0 {
		retVal, err = toml.Marshal(frontMatter)
	}

	return retVal, err
}

//-----------------------------------------------------------------------------
// SerializeFrontMatterToString is a convenience wrapper for string-based callers.
//-----------------------------------------------------------------------------

func SerializeFrontMatterToString(frontMatter map[string]any) (string, error) {

	data, err := SerializeFrontMatter(frontMatter)

	if err != nil {

		return "", err
	}

	return string(bytes.TrimRight(data, "\n")), nil
}

//-----------------------------------------------------------------------------
// DeserializeFrontMatter parses frontmatter content from a serialized block or raw TOML content.
//-----------------------------------------------------------------------------

func DeserializeFrontMatter(data []byte) (map[string]any, error) {

	trimmed := strings.TrimSpace(string(data))

	if trimmed == "" {

		return map[string]any{}, nil
	}

	frontMatter := map[string]any{}

	if _, err := toml.Decode(trimmed, &frontMatter); err != nil {

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
