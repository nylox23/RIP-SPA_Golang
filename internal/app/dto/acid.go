package dto

type AcidCreateRequest struct {
	NameExt   string  `json:"NameExt" binding:"required,max=25"`
	Info      string  `json:"Info" binding:"max=500"`
	Name      string  `json:"Name" binding:"required,max=15"`
	Hplus     int     `json:"Hplus" binding:"required,min=1"`
	MolarMass float32 `json:"MolarMass" binding:"required,min=0"`
}

type AcidUpdateRequest struct {
	NameExt   string  `json:"NameExt" binding:"max=25"`
	Info      string  `json:"Info" binding:"max=500"`
	Name      string  `json:"Name" binding:"max=15"`
	Hplus     int     `json:"Hplus" binding:"min=1"`
	MolarMass float32 `json:"MolarMass" binding:"min=0"`
}

type AcidResponse struct {
	ID        int     `json:"ID"`
	NameExt   string  `json:"NameExt"`
	Info      string  `json:"Info"`
	Name      string  `json:"Name"`
	Hplus     int     `json:"Hplus"`
	MolarMass float32 `json:"MolarMass"`
	Img       string  `json:"Img,omitempty"`
}

type AcidListResponse struct {
	Acids []AcidResponse `json:"Acids"`
}
