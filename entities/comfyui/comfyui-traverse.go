package comfyui

import (
	"encoding/json"
)

type Basic struct {
	Nodes   []Node  `json:"nodes"`
	Version float64 `json:"version"`
}

func (r *Basic) UnmarshalJSON(data []byte) error {
	var err error
	*r, err = UnmarshalIsolatedComfyUI(data)
	if err != nil {
		return err
	}
	return nil
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
	Version float64           `json:"version"`
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
		basic.Nodes = append(basic.Nodes, v)
	}

	if nodeErrors != nil {
		return basic, nodeErrors
	}

	return basic, nil
}
