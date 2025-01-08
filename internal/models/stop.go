package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Stop struct {
	Stop      string
	NameEn    *string
	NameTc    *string
	NameSc    *string
	Latitude  *string
	Longitude *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	CompanyId *int
	Type      *int
	SlugEn    *string
}

type NearestStop struct {
	Stop     string
	Distance float64
}

var stopColumns = []string{"stop", "name_en", "name_tc", "name_sc", "latitude", "longitude", "is_active", "created_at", "updated_at", "company_id", "type", "slug_en"}

func GetNearestStopByStopId(ctx context.Context, dbpool *pgx.Conn, stopId string) (*Stop, error) {
	query := fmt.Sprintf("SELECT %s FROM tbl_stop WHERE is_active = true AND stop = '%s'",
		strings.Join(stopColumns, ","), stopId)
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stops []*Stop
	for rows.Next() {
		var stop Stop

		err := rows.Scan(
			&stop.Stop,
			&stop.NameEn,
			&stop.NameTc,
			&stop.NameSc,
			&stop.Latitude,
			&stop.Longitude,
			&stop.IsActive,
			&stop.CreatedAt,
			&stop.UpdatedAt,
			&stop.CompanyId,
			&stop.Type,
			&stop.SlugEn,
		)
		if err != nil {
			return nil, err
		}

		stops = append(stops, &stop)
	}

	if len(stops) == 0 {
		return nil, nil
	}

	return stops[0], nil
}

func GetActiveNearestStops(ctx context.Context, dbpool *pgx.Conn, transportType int, latitude string, longitude string, limit int) ([]NearestStop, error) {
	query := fmt.Sprintf("SELECT stop, ( 6371000 * acos( cos( radians( %s ) ) * cos( radians(latitude::double precision) ) * cos( radians(longitude::double precision) - radians(%s) ) + sin( radians(%s) ) * sin( radians(latitude::double precision)) ) ) AS distance FROM tbl_stop WHERE is_active = true AND type = %d ORDER BY distance LIMIT %d",
		latitude, longitude, latitude, transportType, limit)
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var nearestStops []NearestStop
	for rows.Next() {
		var nearestStop NearestStop

		err := rows.Scan(
			&nearestStop.Stop,
			&nearestStop.Distance,
		)
		if err != nil {
			return nil, err
		}

		nearestStops = append(nearestStops, nearestStop)
	}

	return nearestStops, nil
}
