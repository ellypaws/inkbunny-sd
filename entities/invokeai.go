package entities

import (
	"encoding/json"
	"errors"
	"fmt"
)

func UnmarshalInvokeAI(data []byte) (InvokeAI, error) {
	var r InvokeAI
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *InvokeAI) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// Convert converts an InvokeAI instance into a TextToImageRequest.
// It maps available fields from InvokeAI to TextToImageRequest.
func (r *InvokeAI) Convert() TextToImageRequest {
	loraHashes := make(map[string]string)
	for _, lora := range r.Loras {
		loraHashes[lora.Model.Hash] = lora.Model.Name
	}

	return TextToImageRequest{
		Prompt:         r.PositivePrompt,
		NegativePrompt: r.NegativePrompt,
		Width:          int(r.Width),
		Height:         int(r.Height),
		Seed:           r.Seed,
		CFGScale:       r.CFGScale,
		Steps:          int(r.Steps),
		SamplerName:    r.Scheduler,
		BatchSize:      1,
		Comments: map[string]string{
			"generation_mode":       r.GenerationMode,
			"app_version":           r.AppVersion,
			"seamless_x":            fmt.Sprintf("%t", r.SeamlessX),
			"seamless_y":            fmt.Sprintf("%t", r.SeamlessY),
			"positive_style_prompt": r.PositiveStylePrompt,
			"negative_style_prompt": r.NegativeStylePrompt,
		},
		OverrideSettings: Config{
			RandnSource:       r.RandDevice,
			SDModelCheckpoint: &r.Model.Name,
			SDCheckpointHash:  r.Model.Hash,
		},
		Scheduler:  &r.Scheduler,
		LoraHashes: loraHashes,
	}
}

type InvokeAI struct {
	GenerationMode       string              `json:"generation_mode"`
	PositivePrompt       string              `json:"positive_prompt"`
	NegativePrompt       string              `json:"negative_prompt"`
	Width                int64               `json:"width"`
	Height               int64               `json:"height"`
	Seed                 int64               `json:"seed"`
	RandDevice           string              `json:"rand_device"`
	CFGScale             float64             `json:"cfg_scale"`
	CFGRescaleMultiplier float64             `json:"cfg_rescale_multiplier"`
	Steps                int64               `json:"steps"`
	Scheduler            string              `json:"scheduler"`
	SeamlessX            bool                `json:"seamless_x"`
	SeamlessY            bool                `json:"seamless_y"`
	Model                InvokeAIModel       `json:"model"`
	Loras                []InvokeAILora      `json:"loras"`
	PositiveStylePrompt  string              `json:"positive_style_prompt"`
	NegativeStylePrompt  string              `json:"negative_style_prompt"`
	Regions              []CanvasEntityState `json:"regions"`
	CanvasV2Metadata     CanvasV2Metadata    `json:"canvas_v2_metadata"`
	AppVersion           string              `json:"app_version"`
}

type InvokeAIModel struct {
	Key  string `json:"key"`
	Hash string `json:"hash"`
	Name string `json:"name"`
	Base string `json:"base"`
	Type string `json:"type"`
}

type InvokeAILora struct {
	Model  InvokeAIModel `json:"model"`
	Weight float64       `json:"weight"`
}

type CanvasV2Metadata struct {
	ReferenceImages  []CanvasReferenceImageState   `json:"referenceImages"`
	ControlLayers    []CanvasControlLayerState     `json:"controlLayers"`
	InpaintMasks     []CanvasInpaintMaskState      `json:"inpaintMasks"`
	RasterLayers     []CanvasRasterLayerState      `json:"rasterLayers"`
	RegionalGuidance []CanvasRegionalGuidanceState `json:"regionalGuidance"`
}

// CanvasEntityState also known as InvokeAI.Regions
type CanvasEntityState interface {
	ID() string
	Type() string
}

type CanvasReferenceImageState struct {
	Id        string               `json:"id"`
	Name      *string              `json:"name"`
	IsEnabled bool                 `json:"isEnabled"`
	IsLocked  bool                 `json:"isLocked"`
	TypeField string               `json:"type"` // "reference_image"
	IPAdapter IPAdapterOrFluxRedux `json:"ipAdapter"`
}

func (c CanvasReferenceImageState) ID() string   { return c.Id }
func (c CanvasReferenceImageState) Type() string { return c.TypeField }

type CanvasControlLayerState struct {
	Id                     string         `json:"id"`
	Name                   *string        `json:"name"`
	IsEnabled              bool           `json:"isEnabled"`
	IsLocked               bool           `json:"isLocked"`
	TypeField              string         `json:"type"` // "control_layer"
	Position               Coordinate     `json:"position"`
	Opacity                float64        `json:"opacity"`
	Objects                CanvasObjects  `json:"objects"`
	WithTransparencyEffect bool           `json:"withTransparencyEffect"`
	ControlAdapter         ControlAdapter `json:"controlAdapter"`
}

type CanvasObjectState interface {
	ID() string
	Type() string
}

type CanvasObjects []CanvasObjectState

func (s *CanvasObjects) UnmarshalJSON(data []byte) error {
	var rawSlice []json.RawMessage
	if err := json.Unmarshal(data, &rawSlice); err != nil {
		return err
	}
	result := make([]CanvasObjectState, len(rawSlice))
	for i, raw := range rawSlice {
		var wrapper canvasObjectStateWrapper
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return err
		}
		result[i] = wrapper.Obj
	}
	*s = result
	return nil
}

// Wrapper for custom unmarshalling of a CanvasObjectState.
type canvasObjectStateWrapper struct {
	Obj CanvasObjectState
}

func (w *canvasObjectStateWrapper) UnmarshalJSON(data []byte) error {
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}
	var typeField string
	if rawType, ok := rawMap["type"]; ok {
		if err := json.Unmarshal(rawType, &typeField); err != nil {
			return err
		}
	} else {
		return errors.New("missing type field in canvas object")
	}

	var err error
	switch typeField {
	case "brush_line":
		w.Obj, err = unmarshalAs[CanvasBrushLineState](data)
	case "brush_line_with_pressure":
		w.Obj, err = unmarshalAs[CanvasBrushLineWithPressureState](data)
	case "eraser_line":
		w.Obj, err = unmarshalAs[CanvasEraserLineState](data)
	case "eraser_line_with_pressure":
		w.Obj, err = unmarshalAs[CanvasEraserLineWithPressureState](data)
	case "rect":
		w.Obj, err = unmarshalAs[CanvasRectState](data)
	case "image":
		w.Obj, err = unmarshalAs[CanvasImageState](data)
	default:
		return fmt.Errorf("unknown canvas object type: %s", typeField)
	}
	return err
}

func unmarshalAs[T CanvasObjectState](data []byte) (T, error) {
	var obj T
	if err := json.Unmarshal(data, &obj); err != nil {
		return obj, err
	}
	return obj, nil
}

func (c CanvasControlLayerState) ID() string   { return c.Id }
func (c CanvasControlLayerState) Type() string { return c.TypeField }

type CanvasInpaintMaskState struct {
	Id        string        `json:"id"`
	Name      *string       `json:"name"`
	IsEnabled bool          `json:"isEnabled"`
	IsLocked  bool          `json:"isLocked"`
	TypeField string        `json:"type"` // "inpaint_mask"
	Position  Coordinate    `json:"position"`
	Fill      Fill          `json:"fill"`
	Opacity   float64       `json:"opacity"`
	Objects   CanvasObjects `json:"objects"`
}

func (c CanvasInpaintMaskState) ID() string   { return c.Id }
func (c CanvasInpaintMaskState) Type() string { return c.TypeField }

type CanvasRasterLayerState struct {
	Id        string        `json:"id"`
	Name      *string       `json:"name"`
	IsEnabled bool          `json:"isEnabled"`
	IsLocked  bool          `json:"isLocked"`
	TypeField string        `json:"type"` // "raster_layer"
	Position  Coordinate    `json:"position"`
	Opacity   float64       `json:"opacity"`
	Objects   CanvasObjects `json:"objects"`
}

func (c CanvasRasterLayerState) ID() string   { return c.Id }
func (c CanvasRasterLayerState) Type() string { return c.TypeField }

type CanvasRegionalGuidanceState struct {
	Id              string                                `json:"id"`
	Name            *string                               `json:"name"`
	IsEnabled       bool                                  `json:"isEnabled"`
	IsLocked        bool                                  `json:"isLocked"`
	TypeField       string                                `json:"type"` // "regional_guidance"
	Position        Coordinate                            `json:"position"`
	Opacity         float64                               `json:"opacity"`
	Objects         CanvasObjects                         `json:"objects"`
	Fill            Fill                                  `json:"fill"`
	PositivePrompt  *string                               `json:"positivePrompt"`
	NegativePrompt  *string                               `json:"negativePrompt"`
	ReferenceImages []RegionalGuidanceReferenceImageState `json:"referenceImages"`
	AutoNegative    bool                                  `json:"autoNegative"`
}

func (c CanvasRegionalGuidanceState) ID() string   { return c.Id }
func (c CanvasRegionalGuidanceState) Type() string { return c.TypeField }

type IPAdapterOrFluxRedux struct {
	Type            string         `json:"type"` // "ip_adapter" or "flux_redux"
	Image           *ImageWithDims `json:"image"`
	Model           *string        `json:"model"`
	Weight          *float64       `json:"weight,omitempty"`
	BeginEndStepPct *[2]float64    `json:"beginEndStepPct,omitempty"`
	Method          *string        `json:"method,omitempty"`
	ClipVisionModel *string        `json:"clipVisionModel,omitempty"`
}

type ControlAdapter struct {
	Type            string      `json:"type"` // "controlnet", "t2i_adapter" or "control_lora"
	Model           *string     `json:"model,omitempty"`
	Weight          *float64    `json:"weight,omitempty"`
	BeginEndStepPct *[2]float64 `json:"beginEndStepPct,omitempty"`
	ControlMode     *string     `json:"controlMode,omitempty"`
}

type RegionalGuidanceReferenceImageState struct {
	Id        string               `json:"id"`
	IPAdapter IPAdapterOrFluxRedux `json:"ipAdapter"`
}

type CanvasImageState struct {
	Id        string        `json:"id"`
	TypeField string        `json:"type"` // "image"
	Image     ImageWithDims `json:"image"`
}

func (c CanvasImageState) ID() string   { return c.Id }
func (c CanvasImageState) Type() string { return c.TypeField }

type CanvasBrushLineState struct {
	Id          string    `json:"id"`
	TypeField   string    `json:"type"` // "brush_line"
	StrokeWidth float64   `json:"strokeWidth"`
	Points      []float64 `json:"points"`
	Color       RgbaColor `json:"color"`
	Clip        *Rect     `json:"clip"`
}

func (c CanvasBrushLineState) ID() string   { return c.Id }
func (c CanvasBrushLineState) Type() string { return c.TypeField }

type CanvasBrushLineWithPressureState struct {
	Id          string    `json:"id"`
	TypeField   string    `json:"type"` // "brush_line_with_pressure"
	StrokeWidth float64   `json:"strokeWidth"`
	Points      []float64 `json:"points"`
	Color       RgbaColor `json:"color"`
	Clip        *Rect     `json:"clip"`
}

func (c CanvasBrushLineWithPressureState) ID() string   { return c.Id }
func (c CanvasBrushLineWithPressureState) Type() string { return c.TypeField }

type CanvasEraserLineState struct {
	Id          string    `json:"id"`
	TypeField   string    `json:"type"` // "eraser_line"
	StrokeWidth float64   `json:"strokeWidth"`
	Points      []float64 `json:"points"`
	Clip        *Rect     `json:"clip"`
}

func (c CanvasEraserLineState) ID() string   { return c.Id }
func (c CanvasEraserLineState) Type() string { return c.TypeField }

type CanvasEraserLineWithPressureState struct {
	Id          string    `json:"id"`
	TypeField   string    `json:"type"` // "eraser_line_with_pressure"
	StrokeWidth float64   `json:"strokeWidth"`
	Points      []float64 `json:"points"`
	Clip        *Rect     `json:"clip"`
}

func (c CanvasEraserLineWithPressureState) ID() string   { return c.Id }
func (c CanvasEraserLineWithPressureState) Type() string { return c.TypeField }

type CanvasRectState struct {
	Id        string    `json:"id"`
	TypeField string    `json:"type"` // "rect"
	Rect      Rect      `json:"rect"`
	Color     RgbaColor `json:"color"`
}

func (c CanvasRectState) ID() string   { return c.Id }
func (c CanvasRectState) Type() string { return c.TypeField }

type Coordinate struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Fill struct {
	Style string   `json:"style"` // e.g. "solid", "grid", etc.
	Color RgbColor `json:"color"`
}

type RgbColor struct {
	R int `json:"r"`
	G int `json:"g"`
	B int `json:"b"`
}

type RgbaColor struct {
	R int     `json:"r"`
	G int     `json:"g"`
	B int     `json:"b"`
	A float64 `json:"a"`
}

type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type ImageWithDims struct {
	ImageName string `json:"image_name"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
