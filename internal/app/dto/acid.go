package dto

type AcidCreateRequest struct {
	NameExt   string  `json:"name_ext" binding:"required,max=25"`
	Info      string  `json:"info" binding:"max=500"`
	Name      string  `json:"name" binding:"required,max=15"`
	Hplus     int     `json:"hplus" binding:"required,min=1"`
	MolarMass float32 `json:"molar_mass" binding:"required,min=0"`
}

type AcidUpdateRequest struct {
	NameExt   string  `json:"name_ext" binding:"max=25"`
	Info      string  `json:"info" binding:"max=500"`
	Name      string  `json:"name" binding:"max=15"`
	Hplus     int     `json:"hplus" binding:"min=1"`
	MolarMass float32 `json:"molar_mass" binding:"min=0"`
}

type AcidResponse struct {
	ID        int     `json:"id"`
	NameExt   string  `json:"name_ext"`
	Info      string  `json:"info"`
	Name      string  `json:"name"`
	Hplus     int     `json:"hplus"`
	MolarMass float32 `json:"molar_mass"`
	Img       string  `json:"img,omitempty"`
}

type AcidListResponse struct {
	Acids []AcidResponse `json:"acids"`
}
