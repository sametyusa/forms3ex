package innsecure

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
	// ErrBadRequest is returned in response to JSON decode errors.
	ErrBadRequest = errors.New("JSON could not be decoded")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
// Useful in a profilesvc server.
func MakeHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	// GET		/hotels/:hotelID/bookings 		retrieves a list of bookings
	// POST		/hotels/:hotelID/bookings 		adds another booking
	// GET		/hotels/:hotelID/bookings/:ID 	adds another booking

	r.Methods("GET").Path("/hotels/{org_id}/bookings").Handler(httptransport.NewServer(
		e.ListBookings,
		httptransport.NopRequestDecoder,
		encodeResponse,
		options...,
	))
	r.Methods("POST").Path("/hotels/{org_id}/bookings").Handler(httptransport.NewServer(
		e.CreateBooking,
		decodeCreateBookingRequest,
		encodeResponseWithStatus(http.StatusCreated),
		options...,
	))
	r.Methods("GET").Path("/hotels/{org_id}/bookings/{id}").Handler(httptransport.NewServer(
		e.GetBookingByID,
		decodeID,
		encodeResponse,
		options...,
	))
	return r
}

func decodeCreateBookingRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var p Booking
	if e := json.NewDecoder(r.Body).Decode(&p); e != nil {
		return nil, ErrBadRequest
	}
	return p, nil
}

func decodeID(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return id, nil
}

// Returns an EncodeResponseFunc that will set the given status code before
// encode the JSON response with encodeResponse.
func encodeResponseWithStatus(code int) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		w.WriteHeader(code)
		return encodeResponse(ctx, w, response)
	}
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	code := codeFrom(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": http.StatusText(code),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrInvalidBooking:
		return http.StatusBadRequest
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case jwt.ErrTokenContextMissing:
		return http.StatusUnauthorized
	case jwt.ErrTokenInvalid:
		return http.StatusUnauthorized
	case jwt.ErrTokenExpired:
		return http.StatusUnauthorized
	case jwt.ErrTokenMalformed:
		return http.StatusUnauthorized
	case jwt.ErrTokenNotActive:
		return http.StatusUnauthorized
	case jwt.ErrUnexpectedSigningMethod:
		return http.StatusUnauthorized
	case jwt.ErrTokenMalformed:
		return http.StatusUnauthorized
	case jwt.ErrTokenExpired:
		return http.StatusUnauthorized
	case jwt.ErrTokenNotActive:
		return http.StatusUnauthorized
	case jwt.ErrTokenInvalid:
		return http.StatusUnauthorized
	case stdjwt.ErrSignatureInvalid:
		return http.StatusUnauthorized
	}
	return http.StatusInternalServerError
}
