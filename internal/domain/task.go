package domain

// Task описывает данные которые будут использоваться для запуска модуля.
type Task struct {
	// Имя функции
	Function string `json:"function" example:"calc"`
	// Список числовых аргументов функции
	Args     []float64 `json:"args"     example:"1.0,2.0,3.0"`
	WasmPath string    `json:"wasmPath" example:"/opt/wapa/runner/data/64b5f827-25b3-44c9-a590-f7fdda440b7c.wasm"`
	JsPath   string    `json:"jsPath"   example:"/opt/wapa/runner/data/64b5f827-25b3-44c9-a590-f7fdda440b7c.js"`
}
