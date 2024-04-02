package sd

import (
	"github.com/ellypaws/inkbunny-sd/entities"
)

func (h *Host) GetConfig() (*entities.Config, error) {
	const configPath = "/sdapi/v1/options"

	body, err := h.GET(configPath)
	if err != nil {
		return nil, err
	}

	var apiConfig entities.Config
	apiConfig, err = entities.UnmarshalConfig(body)
	if err != nil {
		return nil, err
	}

	return &apiConfig, nil
}

func (h *Host) GetCheckpoint() (*string, error) {
	apiConfig, err := h.GetConfig()
	if err != nil {
		return nil, err
	}

	return apiConfig.SDModelCheckpoint, nil
}

func (h *Host) GetVAE() (*string, error) {
	apiConfig, err := h.GetConfig()
	if err != nil {
		return nil, err
	}

	return apiConfig.SDVae, nil
}

func (h *Host) GetHypernetwork() (*string, error) {
	apiConfig, err := h.GetConfig()
	if err != nil {
		return nil, err
	}

	return apiConfig.SDHypernetwork, nil
}
