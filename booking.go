package innsecure

// Booking represents a booking.
type Booking struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Version int    `json:"version"`
	HotelID int    `json:"hotel_id"`
	Arrive  string `json:"arrive"`
	Leave   string `json:"leave"`
	Name    string `json:"name"`
}

// Listing contains a paginated list of bookings.
type Listing struct {
	// Data contains the list of bookings.
	// In a paginated response (not implemented here), this would contain the
	// current page, and additional fields would be added to show the offset
	// and total found.
	Data []Booking `json:"data"`
}
