package domain

// PhoneMetadata holds data about the phone that is reporting the user position.
type PhoneMetadata struct {
	DeviceID    string `json:"device_id"`    // Unique device identifier
	Model       string `json:"model"`        // Model of the phone
	OSVersion   string `json:"os_version"`   // Operating System version
	Carrier     string `json:"carrier"`      // Cellular carrier
	CorporateID string `json:"corporate_id"` // Unique ID assigned by the corporation to the device
	// You can add more fields as required
}
