package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Route struct {
	ID            int
	Type          *int
	Route         string
	RouteSlug     *string
	Bound         *string
	ServiceType   *string
	ServiceNameTc *string
	ServiceNameSc *string
	ServiceNameEn *string
	OriginEn      *string
	OriginTc      *string
	OriginSc      *string
	DestinationEn *string
	DestinationTc *string
	DestinationSc *string
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	CompanyID     *int
	RegionSlug    *string
	RouteID       *string
	RouteEn       *string
	RouteTc       *string
	RouteSc       *string
}

var routeColumns = []string{"id", "type", "route", "route_slug", "bound", "service_type", "service_name_tc", "service_name_sc", "service_name_en", "origin_en", "origin_tc", "origin_sc", "destination_en", "destination_tc", "destination_sc", "is_active", "created_at", "updated_at", "company_id", "region_slug", "route_id", "route_en", "route_tc", "route_sc"}

func GetActiveRoutesByTransportType(ctx context.Context, dbpool *pgx.Conn, transportType int) ([]*Route, error) {
	query := fmt.Sprintf("SELECT %s FROM tbl_route WHERE is_active = true AND type = %d Order By id ASC", strings.Join(routeColumns, ","), transportType)
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*Route
	for rows.Next() {
		var route Route

		err := rows.Scan(
			&route.ID,
			&route.Type,
			&route.Route,
			&route.RouteSlug,
			&route.Bound,
			&route.ServiceType,
			&route.ServiceNameTc,
			&route.ServiceNameSc,
			&route.ServiceNameEn,
			&route.OriginEn,
			&route.OriginTc,
			&route.OriginSc,
			&route.DestinationEn,
			&route.DestinationTc,
			&route.DestinationSc,
			&route.IsActive,
			&route.CreatedAt,
			&route.UpdatedAt,
			&route.CompanyID,
			&route.RegionSlug,
			&route.RouteID,
			&route.RouteEn,
			&route.RouteTc,
			&route.RouteSc,
		)
		if err != nil {
			return nil, err
		}

		routes = append(routes, &route)
	}

	return routes, nil
}
