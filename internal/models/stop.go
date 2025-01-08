package models

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type NearestStop struct {
	Stop     string
	Distance float64
}

func GetActiveNearestStops(ctx context.Context, dbpool *pgx.Conn, latitude string, longitude string, limit int) ([]NearestStop, error) {
	query := fmt.Sprintf("SELECT stop, ( 6371000 * acos( cos( radians( %s ) ) * cos( radians(latitude::double precision) ) * cos( radians(longitude::double precision) - radians(%s) ) + sin( radians(%s) ) * sin( radians(latitude::double precision)) ) ) AS distance FROM tbl_stop WHERE is_active = true",
		latitude, longitude, latitude)
	rows, err := dbpool.Query(ctx, query)
	log.Println(query)
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
