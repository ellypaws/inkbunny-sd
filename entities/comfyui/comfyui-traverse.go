package comfyui

import (
	"encoding/json"
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
		v, err := assertMarshal[Node](node, false)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			continue
		}
		r.Nodes = append(r.Nodes, v)
	}

	return basic, nodeErrors
}
