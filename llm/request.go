// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    request, err := UnmarshalRequest(bytes)
//    bytes, err = request.Marshal()

package llm

import "encoding/json"

func UnmarshalRequest(data []byte) (Request, error) {
	var r Request
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Request) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Role string

const (
	SystemRole Role = "system"
	UserRole   Role = "user"
)

type Request struct {
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
	MaxTokens   int64     `json:"max_tokens"`
	Stream      bool      `json:"stream"`
}

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}
