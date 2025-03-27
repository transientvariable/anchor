package support

import (
	"fmt"

	json "github.com/json-iterator/go"
)

// TODO: Revisit how errors are handled here for the JSON serialization functions. May lead to sensitive data leakage!

// ToJSON serializes the provided source value to a raw JSON byte sequence, which can be cast to a string to get the
// text representation. If an error occurs during serialization, the error message is serialized to JSON, e.g.
// `{"error": "error message..."}`.
func ToJSON(source any) json.RawMessage {
	return toJSON(source, false)
}

// ToJSONFormatted performs the same function as ToJSON and applies formatting to make the result more human-readable.
func ToJSONFormatted(source any) json.RawMessage {
	return toJSON(source, true)
}

func toJSON(source any, indent bool) json.RawMessage {
	var (
		dataJSON []byte
		err      error
	)

	if source != nil {
		if indent {
			dataJSON, err = json.MarshalIndent(&source, "", "  ")
		} else {
			dataJSON, err = json.Marshal(&source)
		}
	}

	if err != nil {
		return json.RawMessage(fmt.Sprintf(`{"error": "%s"}`, err.Error()))
	}
	return dataJSON
}
