package carwise

// @model PredictRequest
// @Description Request body for creating a car price prediction
type PredictRequest struct {
	Brand            string  `json:"Marka" example:"Toyota"`      // Car brand
	Series           string  `json:"Seri" example:"Corolla"`      // Car series
	Model            string  `json:"Model" example:"1.6 Vision"`  // Car model
	Year             int     `json:"Yıl" example:"2019"`          // Manufacturing year
	Mileage          float64 `json:"Kilometre" example:"75000.0"` // Mileage in kilometers
	EngineVolume     float64 `json:"Motor_Hacmi" example:"1.6"`   // Engine volume in liters
	EnginePower      float64 `json:"Motor_Gücü" example:"132.0"`  // Engine power in horsepower
	AccidentHistory  float64 `json:"Tramer" example:"0.0"`        // Number of accidents
	TransmissionType string  `json:"Vites_Tipi" example:"Manuel"` // Transmission type (Manuel/Otomatik)
	FuelType         string  `json:"Yakıt_Tipi" example:"Benzin"` // Fuel type (Benzin/Dizel/LPG)
	BodyType         string  `json:"Kasa_Tipi" example:"Sedan"`   // Body type (Sedan/Hatchback/SUV)
	Color            string  `json:"Renk" example:"Beyaz"`        // Car color
	OriginalParts    int     `json:"Orjinal_sayısı" example:"10"` // Number of original parts
	ReplacedParts    int     `json:"Değişen_sayısı" example:"0"`  // Number of replaced parts
	PaintedParts     int     `json:"Boyalı_sayısı" example:"2"`   // Number of painted parts
	UserId           string  `json:"-"`                           // User ID (from context)
}

// @model GetPredictsRequest
// @Description Request parameters for getting user's predictions
type GetPredictsRequest struct {
	UserId string `json:"-"` // User ID (from context)
	Page   int    `json:"-"` // Page number
	Limit  int    `json:"-"` // Items per page
}

// @model GetPredictsResponse
// @Description Response containing user's predictions
type GetPredictsResponse struct {
	Predicts []Predict `json:"predicts"` // List of predictions
	Total    int       `json:"total"`    // Total number of predictions
}
