package comfyui

import (
	"encoding/json"
	"iter"
)

type Basic struct {
	Nodes   []Node  `json:"nodes"`
	Version float64 `json:"version"`
}

func (r *Basic) UnmarshalJSON(data []byte) error {
	var err error
	*r, err = UnmarshalIsolatedComfyUI(data)
	return err
}

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
	Nodes   []json.RawMessage `json:"nodes"`
	Extra   IsolatedExtra     `json:"extra"`
	Version float64           `json:"version"`
}

type IsolatedExtra struct {
	GroupNodes map[string]json.RawMessage `json:"groupNodes"`
}

func (r *IsolatedComfyUI) parse() (Basic, error) {
	basic := Basic{
		Version: r.Version,
	}

	var nodeErrors NodeErrors
	for node := range r.iterator() {
		v, err := assertMarshal[Node](node, false)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			continue
		}
		basic.Nodes = append(basic.Nodes, v)
	}

	if nodeErrors != nil {
		return basic, nodeErrors
	}

	return basic, nil
}

func (r *IsolatedComfyUI) iterator() iter.Seq[json.RawMessage] {
	return func(yield func(json.RawMessage) bool) {
		for _, node := range r.Nodes {
			if !yield(node) {
				return
			}
		}
		for _, node := range r.Extra.GroupNodes {
			if !yield(node) {
				return
			}
		}
	}
}
