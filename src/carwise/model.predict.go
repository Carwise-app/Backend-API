package carwise

// @model Predict
// @Description Car price prediction model
type Predict struct {
	Id               string  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`         // Unique identifier
	CreatedBy        string  `json:"created_by" example:"550e8400-e29b-41d4-a716-446655440000"` // User who made the prediction
	Brand            string  `json:"brand" example:"Toyota"`                                    // Car brand
	Series           string  `json:"series" example:"Corolla"`                                    // Car series
	Model            string  `json:"model" example:"1.6 Vision"`                                // Car model
	Year             int     `json:"year" example:"2019"`                                        // Manufacturing year
	Mileage          float64 `json:"mileage" example:"75000.0"`                               // Mileage in kilometers
	EngineVolume     float64 `json:"engine_volume" example:"1.6"`                                 // Engine volume in liters
	EnginePower      float64 `json:"engine_power" example:"132.0"`                                // Engine power in horsepower
	AccidentHistory  float64 `json:"accident_history" example:"0.0"`                                      // Number of accidents
	TransmissionType string  `json:"transmission_type" example:"Manuel"`                               // Transmission type (Manuel/Otomatik)
	FuelType         string  `json:"fuel_type" example:"Benzin"`                               // Fuel type (Benzin/Dizel/LPG)
	BodyType         string  `json:"body_type" example:"Sedan"`                                 // Body type (Sedan/Hatchback/SUV)
	Color            string  `json:"color" example:"Beyaz"`                                      // Car color
	OriginalParts    int     `json:"original_parts" example:"10"`                               // Number of original parts
	ReplacedParts    int     `json:"replaced_parts" example:"0"`                                // Number of replaced parts
	PaintedParts     int     `json:"painted_parts" example:"2"`                                 // Number of painted parts
	Price            float64 `json:"price" example:"250000.00"`                                 // Predicted price in TL
	R2               float64 `json:"r2" example:"0.85"`                                         // R2 score of the prediction model
	MAE              float64 `json:"mae" example:"15000.00"`                                    // Mean Absolute Error in TL
	CreatedAt        int64   `json:"created_at" example:"1717238400"`                           // Creation timestamp
}
