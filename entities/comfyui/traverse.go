package comfyui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func assertMarshal[T any](node json.RawMessage, useNumber bool) (T, error) {
	decoder := json.NewDecoder(bytes.NewReader(node))
	if useNumber {
		decoder.UseNumber()
	}
	var v T
	err := decoder.Decode(&v)
	if err != nil {
		return v, NodeError{
			Node: node,
			Err:  ErrInvalidNode,
		}
	}

	return v, nil
}

var ErrInvalidNode = errors.New("invalid node")

type NodeError struct {
	Node json.RawMessage
	Err  error
}

func (e NodeError) Error() string {
	return fmt.Sprintf("cannot parse node: %s, %v", e.Node, e.Err)
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
