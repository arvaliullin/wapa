package domain

// Experiment описывает результаты выполнения эксперимента.
type Experiment struct {
	ID              string           `json:"id"               example:"uid123"`
	DesignID        string           `json:"design_id"        example:"4d9cb632-a0fe-413e-a94b-9353b1e32963"`
	Hostname        string           `json:"hostname"         example:"f87d6b622107"`
	Arch            string           `json:"arch"             example:"amd64"`
	FunctionResults []FunctionResult `json:"function_results"`
}
