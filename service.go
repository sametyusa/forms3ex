package innsecure

import (
	"context"

	"github.com/pborman/uuid"
)

type ErrorString string

// Error satisfies error.
func (e ErrorString) Error() string {
	return string(e)
}

// ErrInvalidBooking is returned if the booking is invalid.
const ErrInvalidBooking ErrorString = "Invalid booking"

// ErrNotFound is used to indicate that a resource has not been found.
// N.B. it is _not_ used when the collection is empty.
const ErrNotFound ErrorString = "Not found"

// ErrDatabase indicates a problem communicating with the database.
const ErrDatabase ErrorString = "Database error"

// ErrUnauthorized indicates that the requested action doesn't match the
// user's permissions set out in the JWT claims
const ErrUnauthorized ErrorString = "Unauthorised"

// Repository represents a collection of bookings in the database.
type Repository interface {
	// Insert creates a new record in the collection, erroring if the ID
	// already exists.
	Insert(ctx context.Context, in Booking) error
	// List returns the contents of the repository.
	List(ctx context.Context, hotelID int) ([]Booking, error)
	// By ID returns a booking by ID.
	ByID(ctx context.Context, hotelID int, ID string) (*Booking, error)
}

// Service provides operations on Bookings.
type Service interface {
	CreateBooking(ctx context.Context, u *User, b Booking) (*Booking, error)
	ListBookings(ctx context.Context, u *User) (listing *Listing, err error)
	GetBookingByID(ctx context.Context, u *User, ID string) (*Booking, error)
}

type User struct {
	Name    string
	Admin   bool
	HotelID int
}

// NewBookingService returns a pointer to a new booking service instance.
func NewBookingService(r Repository) *BookingService {
	return &BookingService{
		r: r,
	}
}

// BookingService satisfies Service.
type BookingService struct {
	r Repository
}

// ListBookings returns a list of bookings from the database.
func (svc *BookingService) ListBookings(ctx context.Context, u *User) (*Listing, error) {
	if u == nil {
		return nil, ErrUnauthorized
	}

	list, err := svc.r.List(ctx, u.HotelID)
	if err != nil {
		return nil, err
	}

	return &Listing{
		Data: list,
	}, nil
}

func (svc *BookingService) bookingIsValid(b Booking, allowID bool) bool {
	if b.Type != "Booking" || b.Version != 0 || b.HotelID == 0 {
		return false
	}
	if !allowID && b.ID != "" {
		return false
	}
	return true
}

// CreateBooking adds a booking to the collection. It returns a a booking
// object, updated to include its generated ID.
func (svc *BookingService) CreateBooking(ctx context.Context, u *User, b Booking) (*Booking, error) {
	if !svc.bookingIsValid(b, false) {
		return nil, ErrInvalidBooking
	}

	if u == nil {
		return nil, ErrUnauthorized
	}

	if u == nil || !u.Admin || b.HotelID != u.HotelID {
		return nil, ErrUnauthorized
	}

	b.ID = uuid.New()

	err := svc.r.Insert(ctx, b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// GetBookingByID retrieves a booking matching the given ID, if present.
func (svc *BookingService) GetBookingByID(ctx context.Context, u *User, ID string) (*Booking, error) {
	if u == nil {
		return nil, ErrUnauthorized
	}

	b, err := svc.r.ByID(ctx, u.HotelID, ID)
	if err != nil {
		return nil, ErrDatabase
	}

	if b == nil {
		return nil, ErrNotFound
	}

	return b, nil
}

// convertDBError converts a repository error to a domain one.
func convertDBError(err error) error {
	switch err {
	case nil:
		return nil
	case ErrNotFound:
		return ErrNotFound
	default:
		return ErrDatabase
	}
}
