package innsecure_test

import (
	"context"
	"errors"
	"testing"

	"github.com/form3tech/innsecure"
	"github.com/go-kit/kit/endpoint"
)

func noopMiddleware(e endpoint.Endpoint) endpoint.Endpoint {
	return e
}

type svc struct {
	listBookings   func(ctx context.Context, u *innsecure.User) (listing *innsecure.Listing, err error)
	createBooking  func(ctx context.Context, p innsecure.Booking) (*innsecure.Booking, error)
	getBookingByID func(ctx context.Context, u *innsecure.User, ID string) (*innsecure.Booking, error)
}

func (s svc) ListBookings(ctx context.Context, u *innsecure.User) (listing *innsecure.Listing, err error) {
	return s.listBookings(ctx, u)
}
func (s svc) CreateBooking(ctx context.Context, u *innsecure.User, p innsecure.Booking) (*innsecure.Booking, error) {
	return s.createBooking(ctx, p)
}
func (s svc) GetBookingByID(ctx context.Context, u *innsecure.User, ID string) (*innsecure.Booking, error) {
	return s.getBookingByID(ctx, u, ID)
}

func TestCanWrapList(t *testing.T) {
	want := &innsecure.Listing{}
	wantErr := errors.New("testerr")
	service := svc{
		listBookings: func(_ context.Context, _ *innsecure.User) (*innsecure.Listing, error) {
			return want, wantErr
		},
	}

	sut := innsecure.MakeServerEndpoints(service, noopMiddleware)
	got, err := sut.ListBookings(context.TODO(), nil)
	if got != want {
		t.Fatalf("want=%+v, got=%+v", want, got)
	}
	if err != wantErr {
		t.Fatalf("want=%s, got=%s", wantErr, err)
	}
}

func TestCanWrapCreate(t *testing.T) {
	wantIn := innsecure.Booking{ID: "INPUT"}
	wantOut := &innsecure.Booking{ID: "OUTPUT"}
	wantErr := errors.New("testerr")
	service := svc{
		createBooking: func(_ context.Context, in innsecure.Booking) (*innsecure.Booking, error) {
			if in.ID != wantIn.ID {
				t.Fatalf("unexpected input")
			}
			return wantOut, wantErr
		},
	}

	sut := innsecure.MakeServerEndpoints(service, noopMiddleware)
	got, err := sut.CreateBooking(context.TODO(), wantIn)
	if got != wantOut {
		t.Fatalf("want=%+v, got=%+v", wantOut, got)
	}
	if err != wantErr {
		t.Fatalf("want=%s, got=%s", wantErr, err)
	}
}

func TestCanWrapGetByID(t *testing.T) {
	wantIn := "INPUT"
	wantOut := &innsecure.Booking{ID: "OUTPUT"}
	wantErr := errors.New("testerr")
	service := svc{
		getBookingByID: func(_ context.Context, _ *innsecure.User, in string) (*innsecure.Booking, error) {
			if in != wantIn {
				t.Fatalf("want=%s, got=%s", wantIn, in)
			}
			return wantOut, wantErr
		},
	}

	sut := innsecure.MakeServerEndpoints(service, noopMiddleware)
	got, err := sut.GetBookingByID(context.TODO(), wantIn)
	if got != wantOut {
		t.Fatalf("want=%+v, got=%+v", wantOut, got)
	}
	if err != wantErr {
		t.Fatalf("want=%s, got=%s", wantErr, err)
	}
}
