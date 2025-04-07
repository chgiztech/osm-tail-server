package dto

type CoordinateDto struct {
	MinLat float64 `json:"minLat" validate:"required,min=0"`
	MinLon float64 `json:"minLon" validate:"required,min=0"`
	MaxLat float64 `json:"maxLat" validate:"required,min=0"`
	MaxLon float64 `json:"maxLon" validate:"required,min=0"`
}
