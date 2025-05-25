package domain

// Design описывает план эксперимента
type Design struct {
	// Идентификатор плана эксперимента
	ID string `json:"id" example:"exp123"`
	// Имя эксперимента
	Name string `json:"name" example:"x2Integrate"`
	// Язык реализации модуля
	Lang string `json:"lang" example:"go"`
	// Имя файла JavaScript
	JS string `json:"js" example:"index.js"`
	// Имя модуля WebAssembly
	Wasm string `json:"wasm,omitempty" example:"module.wasm"`
	// Массив функций и их аргументов
	Functions []Function `json:"functions"`
}

// Function описывает структуру массива функций с именем функции и списком аргументов
type Function struct {
	// Имя функции
	Function string `json:"function" example:"calc"`
	// Список числовых аргументов функции
	Args []float64 `json:"args" example:"1.0,2.0,3.0"`
}
