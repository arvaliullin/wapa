package domain

type DesignPayload struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Lang      string     `json:"lang"`
	JS        string     `json:"js"`
	Wasm      string     `json:"wasm"`
	Repeats   int        `json:"repeats"`
	Warmup    bool       `json:"warmup"`
	Functions []Function `json:"functions"`
}
