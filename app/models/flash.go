package models

type DataFrameFlash struct {
	Temperature Temperature `json:"temperature,omitempty"`
	PowerDissipation PowerDissipation `json:"power_dissipation,omitempty"`
}

type Temperature struct {
	Values []float64 `json:"values,omitempty"`
}

type PowerDissipation struct {
	Values []float64 `json:"values,omitempty"`
}
