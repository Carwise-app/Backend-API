package carwise

type Predict struct {
	Id               string  `json:"id"`
	CreatedBy        string  `json:"created_by"`
	CreatedAt        int64   `json:"created_at"`
	Brand            string  `json:"brand"`
	Series           string  `json:"series"`
	Model            string  `json:"model"`
	Year             int     `json:"year"`
	Mileage          float64 `json:"mileage"`
	EngineVolume     float64 `json:"engine_volume"`
	EnginePower      float64 `json:"engine_power"`
	AccidentHistory  float64 `json:"accident_history"`
	TransmissionType string  `json:"transmission_type"`
	FuelType         string  `json:"fuel_type"`
	BodyType         string  `json:"body_type"`
	Color            string  `json:"color"`
	OriginalParts    int     `json:"original_parts"`
	ReplacedParts    int     `json:"replaced_parts"`
	PaintedParts     int     `json:"painted_parts"`
	Price            float64 `json:"price"`
	R2               float64 `json:"r2"`
	MAE              float64 `json:"mae"`
}
