package carwise

type PredictRequest struct {
	Brand            string  `json:"Marka"`
	Series           string  `json:"Seri"`
	Model            string  `json:"Model"`
	Year             int     `json:"Yıl"`
	Mileage          float64 `json:"Kilometre"`
	EngineVolume     float64 `json:"Motor_Hacmi"`
	EnginePower      float64 `json:"Motor_Gücü"`
	AccidentHistory  float64 `json:"Tramer"`
	TransmissionType string  `json:"Vites_Tipi"`
	FuelType         string  `json:"Yakıt_Tipi"`
	BodyType         string  `json:"Kasa_Tipi"`
	Color            string  `json:"Renk"`
	OriginalParts    int     `json:"Orjinal_sayısı"`
	ReplacedParts    int     `json:"Değişen_sayısı"`
	PaintedParts     int     `json:"Boyalı_sayısı"`
	UserId           string  `json:"-"`
}

type GetPredictsRequest struct {
	UserId string
	Page   int
	Limit  int
}

type GetPredictsResponse struct {
	Predicts []Predict `json:"predicts"`
	Total    int       `json:"total"`
}
