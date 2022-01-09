package innsecure

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

const UserContextKey = "user"

// Endpoints collects all of the service's endpoints.
type Endpoints struct {
	ListBookings   endpoint.Endpoint
	CreateBooking  endpoint.Endpoint
	GetBookingByID endpoint.Endpoint
}

func contextToUser(ctx context.Context) *User {
	u, ok := ctx.Value(UserContextKey).(*User)
	if !ok {
		return nil
	}
	return u
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service.
func MakeServerEndpoints(s Service, jwtmw endpoint.Middleware) Endpoints {

	return Endpoints{
		ListBookings:   jwtmw(MakeListBookingsEndpoint(s)),
		CreateBooking:  jwtmw(MakeCreateBookingEndpoint(s)),
		GetBookingByID: jwtmw(MakeGetBookingByIDEndpoint(s)),
	}
}

// MakeListBookingsEndpoint returns an endpoint wrapping the given server.
func MakeListBookingsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (response interface{}, err error) {
		u := contextToUser(ctx)
		return s.ListBookings(ctx, u)
	}
}

// MakeCreateBookingEndpoint returns an endpoint wrapping the given server.
func MakeCreateBookingEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		b, ok := request.(Booking)
		if !ok {
			return nil, errors.New("invalid request type, likely bad wiring")
		}
		u := contextToUser(ctx)

		return s.CreateBooking(ctx, u, b)
	}
}

// MakeGetBookingByIDEndpoint returns an endpoint wrapping the given server.
func MakeGetBookingByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id, ok := request.(string)
		if !ok {
			return nil, errors.New("invalid request type, likely bad wiring")
		}
		u := contextToUser(ctx)
		return s.GetBookingByID(ctx, u, id)
	}
}
