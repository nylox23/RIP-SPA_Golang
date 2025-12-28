package dto

type CalcRequestItem struct {
	ID        uint    `json:"id"`
	CMass     float32 `json:"c_mass"`
	AMass     float32 `json:"a_mass"`
	HPlus     int     `json:"h_plus"`
	MolarMass float32 `json:"molar"`
}

type CalcRequest struct {
	Items []CalcRequestItem `json:"items"`
}

type CalcResultItem struct {
	ID     uint    `json:"id"`
	Result float32 `json:"result"`
}

type CalcCallbackRequest struct {
	Results []CalcResultItem `json:"results"`
}
