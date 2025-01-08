package models

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type Route struct {
	ID   int
	Type *int
	// Route         string
	// RouteSlug     *string
	// Bound         *string
	// ServiceType   *string
	// ServiceNameTc *string
	// ServiceNameSc *string
	// ServiceNameEn *string
	// OriginEn      *string
	// OriginTc      *string
	// OriginSc      *string
	// DestinationEn *string
	// DestinationTc *string
	// DestinationSc *string
	// IsActive      bool
	// CreatedAt     time.Time
	// UpdatedAt     time.Time
	// CompanyID     *int
	// RegionSlug    *string
	// RouteID       *string
	// RouteEn       *string
	// RouteTc       *string
	// RouteSc       *string
}

func GetAllRoutes(ctx context.Context, dbpool *pgx.Conn) ([]*Route, error) {
	query := "SELECT ID, type FROM tbl_route LIMIT 500"
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var route Route

		err := rows.Scan(
			&route.ID,
			&route.Type,
		)
		if err != nil {
			return nil, err
		}

		if route.Type != nil {
			log.Printf("Route: {ID: %d, Type: %d}", route.ID, *route.Type)
		} else {
			log.Printf("Route: {ID: %d, Type: nil}", route.ID)
		}
	}

	return nil, nil
}

// func ReadAllCategories(ctx context.Context) ([]*Category, error) {
// 	rows, err := tx.Query(ctx, fmt.Sprintf("SELECT %s FROM categories ORDER BY position LIMIT 500", strings.Join(categoryColumns, ",")))
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
// }
