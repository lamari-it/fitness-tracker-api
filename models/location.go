package models

// Location is an embeddable struct for storing location data
// Used by User (personal location) and TrainerProfile (business location)
type Location struct {
	Latitude    *float64 `gorm:"type:decimal(10,8)" json:"latitude,omitempty"`
	Longitude   *float64 `gorm:"type:decimal(11,8)" json:"longitude,omitempty"`
	CountryCode *string  `gorm:"type:varchar(2)" json:"country_code,omitempty"`
	Region      *string  `gorm:"type:varchar(100)" json:"region,omitempty"`
	City        *string  `gorm:"type:varchar(100)" json:"city,omitempty"`
	District    *string  `gorm:"type:varchar(100)" json:"district,omitempty"`
	PostalCode  *string  `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	RawAddress  *string  `gorm:"type:text" json:"raw_address,omitempty"`
}

// LocationResponse is the DTO for location data in API responses
type LocationResponse struct {
	Latitude    *float64 `json:"latitude,omitempty"`
	Longitude   *float64 `json:"longitude,omitempty"`
	CountryCode *string  `json:"country_code,omitempty"`
	Region      *string  `json:"region,omitempty"`
	City        *string  `json:"city,omitempty"`
	District    *string  `json:"district,omitempty"`
	PostalCode  *string  `json:"postal_code,omitempty"`
	RawAddress  *string  `json:"raw_address,omitempty"`
}

// LocationUpdateRequest is the DTO for updating location data
type LocationUpdateRequest struct {
	Latitude    *float64 `json:"latitude,omitempty" binding:"omitempty,min=-90,max=90"`
	Longitude   *float64 `json:"longitude,omitempty" binding:"omitempty,min=-180,max=180"`
	CountryCode *string  `json:"country_code,omitempty" binding:"omitempty,len=2"`
	Region      *string  `json:"region,omitempty" binding:"omitempty,max=100"`
	City        *string  `json:"city,omitempty" binding:"omitempty,max=100"`
	District    *string  `json:"district,omitempty" binding:"omitempty,max=100"`
	PostalCode  *string  `json:"postal_code,omitempty" binding:"omitempty,max=20"`
	RawAddress  *string  `json:"raw_address,omitempty" binding:"omitempty,max=500"`
}

// HasCoordinates returns true if both latitude and longitude are set
func (l *Location) HasCoordinates() bool {
	return l.Latitude != nil && l.Longitude != nil
}

// HasStructuredLocation returns true if any structured location field is set
func (l *Location) HasStructuredLocation() bool {
	return l.CountryCode != nil || l.Region != nil || l.City != nil || l.District != nil || l.PostalCode != nil
}

// HasAnyLocation returns true if any location data is set
func (l *Location) HasAnyLocation() bool {
	return l.HasCoordinates() || l.HasStructuredLocation() || l.RawAddress != nil
}

// ToResponse converts Location to LocationResponse
func (l *Location) ToResponse() *LocationResponse {
	if !l.HasAnyLocation() {
		return nil
	}
	return &LocationResponse{
		Latitude:    l.Latitude,
		Longitude:   l.Longitude,
		CountryCode: l.CountryCode,
		Region:      l.Region,
		City:        l.City,
		District:    l.District,
		PostalCode:  l.PostalCode,
		RawAddress:  l.RawAddress,
	}
}

// UpdateFromRequest updates Location fields from LocationUpdateRequest
func (l *Location) UpdateFromRequest(req *LocationUpdateRequest) {
	if req == nil {
		return
	}
	if req.Latitude != nil {
		l.Latitude = req.Latitude
	}
	if req.Longitude != nil {
		l.Longitude = req.Longitude
	}
	if req.CountryCode != nil {
		l.CountryCode = req.CountryCode
	}
	if req.Region != nil {
		l.Region = req.Region
	}
	if req.City != nil {
		l.City = req.City
	}
	if req.District != nil {
		l.District = req.District
	}
	if req.PostalCode != nil {
		l.PostalCode = req.PostalCode
	}
	if req.RawAddress != nil {
		l.RawAddress = req.RawAddress
	}
}
