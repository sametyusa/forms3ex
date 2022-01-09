package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/form3tech/innsecure"
)

type BookingRepo struct {
	db *sql.DB
}

// NewRepo returns a new repository backed by the given DB.
func NewRepo(db *sql.DB) *BookingRepo {
	return &BookingRepo{
		db: db,
	}
}

// Insert satisfies Repository.
func (r *BookingRepo) Insert(ctx context.Context, b innsecure.Booking) error {
	stmt := fmt.Sprintf(`insert into "Bookings" ("id", "hotelid", "arrive", "leave", "name") values ('%s', '%d', '%s', '%s', '%s')`, b.ID, b.HotelID, b.Arrive, b.Leave, b.Name)
	_, err := r.db.Exec(stmt)
	return err
}

// List returns the full contents of the repository.
func (r *BookingRepo) List(ctx context.Context, hotelID int) ([]innsecure.Booking, error) {
	stmt := fmt.Sprintf(`select "id", "hotelid", "arrive", "leave", "name" from "Bookings" WHERE hotelid='%d'`, hotelID)
	rows, err := r.db.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to list bookings: %w", err)
	}
	result := []innsecure.Booking{}
	for rows.Next() {
		var b innsecure.Booking
		err = rows.Scan(&b.ID, &b.HotelID, &b.Arrive, &b.Leave, &b.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		result = append(result, b)
	}
	return result, nil
}

// ByID returns a single booking by ID.
// If no booking is found with the given ID, no error is returned.
func (r *BookingRepo) ByID(ctx context.Context, hotelID int, ID string) (*innsecure.Booking, error) {
	q := fmt.Sprintf(`select "id", "hotelid", "arrive", "leave", "name" from "Bookings" where "hotelid"='%d' AND "id"='%s'`, hotelID, ID)
	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b innsecure.Booking
		err = rows.Scan(&b.ID, &b.HotelID, &b.Arrive, &b.Leave, &b.Name)
		if err == nil {
			return &b, nil
		}
	}
	return nil, nil
}
