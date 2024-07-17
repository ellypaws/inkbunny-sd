package comfyui

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type PreNode any

func UnmarshalIsolatedComfyUI(data []byte) (Basic, error) {
	var r IsolatedComfyUI
	err := json.Unmarshal(data, &r)
	if err != nil {
		return Basic{}, err
	}
	return r.parse()
}

func (r *IsolatedComfyUI) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type IsolatedComfyUI struct {
	Nodes   []any   `json:"nodes"`
	Version float64 `json:"version"`
}

func (r *IsolatedComfyUI) parse() (Basic, error) {
	basic := Basic{
		Version: r.Version,
	}
	var nodeErrors NodeErrors

	for _, node := range r.Nodes {
		bin, err := json.Marshal(node)
		if err != nil {
			nodeErrors = append(nodeErrors,
				NodeError{
					Node: node,
					Err:  err,
				})
			continue
		}

		var n Node
		err = json.Unmarshal(bin, &n)
		if err != nil {
			nodeErrors = append(nodeErrors,
				NodeError{
					Node: node,
					Err:  ErrInvalidNode,
					bin:  bin,
				})
			continue
		}
		basic.Nodes = append(basic.Nodes, n)
	}

	if nodeErrors != nil {
		return basic, nodeErrors
	}

	return basic, nil
}

var ErrInvalidNode = errors.New("invalid node")

type NodeError struct {
	Node PreNode
	Err  error
	bin  []byte
}

func (e NodeError) Error() string {
	if e.bin != nil {
		return fmt.Sprintf("cannot parse node: %s, %v", e.bin, e.Err)
	}
	return fmt.Sprintf("cannot parse node: %v, %v", e.Node, e.Err)
}

func (e NodeError) Unwrap() error {
	return e.Err
}

type NodeErrors []error

func (e NodeErrors) Error() string {
	var s strings.Builder
	for i, err := range e {
		if i > 0 {
			s.WriteString("\n")
		}
		s.WriteString(err.Error())
	}
	return s.String()
}

func (e NodeErrors) Unwrap() []error {
	if len(e) == 0 {
		return nil
	}
	return e
}

func (e NodeErrors) Strings() []string {
	var s []string
	for _, err := range e {
		s = append(s, err.Error())
	}
	return s
}

func (e NodeErrors) Push(err error) NodeErrors {
	return append(e, err)
}

func (e NodeErrors) Len() int {
	return len(e)
}
