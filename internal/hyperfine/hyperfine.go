package hyperfine

// HyperfineResult результат измерения утилитой hyperfine
type HyperfineResult struct {
	Results []struct {
		Command   string    `json:"command"`
		Mean      float64   `json:"mean"`
		Stddev    float64   `json:"stddev"`
		Median    float64   `json:"median"`
		User      float64   `json:"user"`
		System    float64   `json:"system"`
		Min       float64   `json:"min"`
		Max       float64   `json:"max"`
		Times     []float64 `json:"times"`
		ExitCodes []int     `json:"exit_codes"`
	} `json:"results"`
}
