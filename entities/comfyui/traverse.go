package comfyui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func assertMarshal[T any](node any, useNumber bool) (T, error) {
	var v T
	bin, err := json.Marshal(node)
	if err != nil {
		return v, NodeError{
			Node: node,
			Err:  err,
		}
	}

	decoder := json.NewDecoder(bytes.NewReader(bin))
	if useNumber {
		decoder.UseNumber()
	}
	err = decoder.Decode(&v)
	if err != nil {
		return v, NodeError{
			Node: node,
			Err:  ErrInvalidNode,
			bin:  bin,
		}

	}

	return v, nil
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
