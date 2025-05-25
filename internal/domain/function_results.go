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
