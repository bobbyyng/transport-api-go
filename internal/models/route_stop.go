package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type RouteStop struct {
	ID          int
	Route       string
	Bound       string
	ServiceType string
	Sequence    int
	Stop        string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompanyId   *int
	Type        *int
	RegionSlug  *string
}

var routeStopColumns = []string{"id", "route", "bound", "service_type", "sequence", "stop", "is_active", "created_at", "updated_at", "company_id", "type", "region_slug"}

func GetActiveRouteStopsByStopId(ctx context.Context, dbpool *pgx.Conn, transportType int, stopId string, serviceType string) ([]*RouteStop, error) {
	query := fmt.Sprintf("SELECT %s FROM tbl_route_stop WHERE is_active = true AND type = %d AND stop = '%s' AND service_type = '%s' Order By id ASC",
		strings.Join(routeStopColumns, ","), transportType, stopId, serviceType)
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routeStops []*RouteStop
	for rows.Next() {
		var routeStop RouteStop

		err := rows.Scan(
			&routeStop.ID,
			&routeStop.Route,
			&routeStop.Bound,
			&routeStop.ServiceType,
			&routeStop.Sequence,
			&routeStop.Stop,
			&routeStop.IsActive,
			&routeStop.CreatedAt,
			&routeStop.UpdatedAt,
			&routeStop.CompanyId,
			&routeStop.Type,
			&routeStop.RegionSlug,
		)
		if err != nil {
			return nil, err
		}
		routeStops = append(routeStops, &routeStop)
	}

	return routeStops, nil
}
