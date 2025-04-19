package domain

// FunctionResult описывает результаты выполнения функции
type FunctionResult struct {
	ID           string    `json:"id" example:"uid123"`
	ExperimentID string    `json:"experiment_id" example:"4d9cb632-a0fe-413e-a94b-9353b1e32963"`
	FunctionName string    `json:"function_name" example:"x2Integrate"`
	Args         []float64 `json:"args" example:"1.0,2.0,3.0"`
	Repeats      int       `json:"repeats" example:"9999"`
	Result       []string  `json:"result" example:"333233.3449995302"`
	Metrics      Metrics   `json:"metrics"`
}

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
