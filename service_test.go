package innsecure_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/form3tech/innsecure"
	"github.com/pborman/uuid"
)

// repo is a basic Repository mock.
type repo struct {
	insert func(ctx context.Context, in innsecure.Booking) error
	list   func(ctx context.Context, hotelID int) ([]innsecure.Booking, error)
	byID   func(ctx context.Context, hotelID int, ID string) (*innsecure.Booking, error)
}

func (r repo) Insert(ctx context.Context, p innsecure.Booking) error {
	return r.insert(ctx, p)
}

func (r repo) List(ctx context.Context, hotelID int) ([]innsecure.Booking, error) {
	return r.list(ctx, hotelID)
}

func (r repo) ByID(ctx context.Context, hotelID int, ID string) (*innsecure.Booking, error) {
	return r.byID(ctx, hotelID, ID)
}

func normalUser() *innsecure.User {
	return &innsecure.User{
		Name:    "Geoff Capes",
		HotelID: 123,
		Admin:   false,
	}
}

func adminUser() *innsecure.User {
	return &innsecure.User{
		Name:    "Geoff Capes",
		HotelID: 123,
		Admin:   true,
	}
}

func TestBookingServiceImplementsService(t *testing.T) {
	var sut interface{} = &innsecure.BookingService{}
	_, ok := sut.(innsecure.Service)
	if !ok {
		t.Fatal("BookingService does not satisfy Service.")
	}
}

func validBooking(ID string) innsecure.Booking {
	return innsecure.Booking{
		ID:      ID,
		Type:    "Booking",
		Version: 0,
		HotelID: 123,
		Arrive:  "2021-08-13",
		Leave:   "2021-08-15",
		Name:    "Jane Guest",
	}
}

// Get bookings list

func TestCanGetEmptyBookingsList(t *testing.T) {
	r := repo{
		list: func(_ context.Context, _ int) ([]innsecure.Booking, error) {
			return []innsecure.Booking{}, nil
		},
	}
	sut := innsecure.NewBookingService(r)
	got, err := sut.ListBookings(context.TODO(), normalUser())

	if err != nil {
		t.Fatal(err)
	}

	if len(got.Data) != 0 {
		t.Fatalf("expected empty list, got %+v", got.Data)
	}
}

func TestCanGetBookingsList(t *testing.T) {
	r := repo{
		list: func(_ context.Context, _ int) ([]innsecure.Booking, error) {
			return []innsecure.Booking{
				{ID: "A"},
				{ID: "B"},
				{ID: "C"},
			}, nil
		},
	}
	sut := innsecure.NewBookingService(r)
	got, err := sut.ListBookings(context.TODO(), normalUser())
	if err != nil {
		t.Fatal(err)
	}

	want := &innsecure.Listing{
		Data: []innsecure.Booking{
			{ID: "A"},
			{ID: "B"},
			{ID: "C"},
		},
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want=%+v, got=%+v", want, got)
	}
}

func TestCanHandleRepoFailureWhenGettingBookingList(t *testing.T) {
	r := repo{
		list: func(_ context.Context, _ int) ([]innsecure.Booking, error) {
			return nil, errors.New("test error")
		},
	}
	sut := innsecure.NewBookingService(r)
	_, err := sut.ListBookings(context.TODO(), normalUser())
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestCanCreateBooking(t *testing.T) {
	var inserted innsecure.Booking
	r := repo{
		insert: func(_ context.Context, in innsecure.Booking) error {
			inserted = in
			return nil
		},
	}
	sut := innsecure.NewBookingService(r)
	b := validBooking("")
	b.HotelID = 123
	got, err := sut.CreateBooking(context.TODO(), adminUser(), b)
	if err != nil {
		t.Fatal(err)
	}

	// Expect ID to be generated.
	gotID := uuid.Parse(got.ID)
	if gotID == nil {
		t.Fatalf("ID ('%s') cannot be parsed as UUID.", got.ID)
	}

	if inserted.ID != got.ID {
		t.Fatalf("want=%s, got=%s", got.ID, inserted.ID)
	}
	if inserted.HotelID != b.HotelID {
		t.Fatalf("want=%d, got=%d", b.HotelID, inserted.HotelID)
	}
}

func TestCanRejectCreateWithNonAdminUser(t *testing.T) {
	r := repo{
		insert: func(_ context.Context, in innsecure.Booking) error {
			return nil
		},
	}
	sut := innsecure.NewBookingService(r)
	b := validBooking("")
	b.HotelID = 123
	_, err := sut.CreateBooking(context.TODO(), normalUser(), b)
	if err != innsecure.ErrUnauthorized {
		t.Fatal("Expected unauthorised error, got none")
	}
}

func TestCanRejectCreateWithInvalidBooking(t *testing.T) {
	cases := map[string]innsecure.Booking{
		"Wrong type": {
			Type:    "Debit",
			Version: 0,
			HotelID: 234,
			Arrive:  "2021-08-13",
			Leave:   "2021-08-15",
			Name:    "Jane Guest",
		},
		"ID in create": {
			Type:    "Booking",
			ID:      "a82a8dc8-a044-4769-970e-2143d9a1a050",
			Version: 0,
			HotelID: 345,
			Arrive:  "2021-08-13",
			Leave:   "2021-08-15",
			Name:    "Jane Guest",
		},
		"Wrong version": {
			Type:    "Booking",
			Version: 1,
			HotelID: 456,
			Arrive:  "2021-08-13",
			Leave:   "2021-08-15",
			Name:    "Jane Guest",
		},
		"Missing hotelID": {
			Type:    "Booking",
			Version: 0,
			Arrive:  "2021-08-13",
			Leave:   "2021-08-15",
			Name:    "Jane Guest",
		},
	}
	for k, b := range cases {
		t.Run(k, func(t *testing.T) {
			r := repo{
				insert: func(_ context.Context, _ innsecure.Booking) error {
					t.Fatal("insert should not have been called, was")
					return errors.New("test error")
				},
			}
			sut := innsecure.NewBookingService(r)
			_, err := sut.CreateBooking(context.TODO(), adminUser(), b)
			if err != innsecure.ErrInvalidBooking {
				t.Fatalf("want=%s, got=%s", innsecure.ErrInvalidBooking, err)
			}
		})
	}
}

func TestCanHandleRepoFailureWhenCreatingBooking(t *testing.T) {
	r := repo{
		insert: func(_ context.Context, in innsecure.Booking) error {
			t.Fatal("record should not have been inserted, was")
			return errors.New("test error")
		},
	}
	sut := innsecure.NewBookingService(r)
	p := innsecure.Booking{}
	_, err := sut.CreateBooking(context.TODO(), adminUser(), p)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

// Get booking

func TestCanGetBookingsByID(t *testing.T) {
	r := repo{
		byID: func(_ context.Context, hotelID int, ID string) (*innsecure.Booking, error) {
			switch ID {
			case "found":
				return &innsecure.Booking{HotelID: 123}, nil
			case "notfound":
				return nil, nil
			case "dberror":
				return nil, errors.New("test error")
			default:
				t.Fatal("unexpected ID")
				return nil, nil
			}
		},
	}
	sut := innsecure.NewBookingService(r)

	cases := []struct {
		id      string
		want    *innsecure.Booking
		wantErr error
	}{
		{id: "found", want: &innsecure.Booking{HotelID: 123}},
		{id: "notfound", want: nil, wantErr: innsecure.ErrNotFound},
		{id: "dberror", want: nil, wantErr: innsecure.ErrDatabase},
	}
	for _, c := range cases {
		t.Run(c.id, func(t *testing.T) {
			got, err := sut.GetBookingByID(context.TODO(), normalUser(), c.id)
			if c.wantErr != nil {
				if err != c.wantErr {
					t.Fatalf("error: want=%+v, got=%+v", c.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(c.want, got) {
				t.Fatalf("want=%+v, got=%+v", c.want, got)
			}
		})
	}
}
