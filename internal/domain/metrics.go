package domain

// Metrics описывает статистические метрики для времени выполнения
type Metrics struct {
	Mean   float64 `json:"mean" example:"0.022032815420000004"`
	Stddev float64 `json:"stddev" example:"0.0005590838760400804"`
	Median float64 `json:"median" example:"0.03726900000037858"`
	User   float64 `json:"user" example:"0.03726900000037858"`
	System float64 `json:"system" example:"0.03726900000037858"`
	Min    float64 `json:"min" example:"0.03685800000312156"`
	Max    float64 `json:"max" example:"2.5455180000026303"`
}
