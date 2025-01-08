package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type Company struct {
	ID        int
	SlugEn    *string
	SlugTc    *string
	SlugSc    *string
	NameEn    *string
	NameTc    *string
	NameSc    *string
	BgColor   *string
	TextColor *string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

var companyColumns = []string{"id", "slug_en", "slug_tc", "slug_sc", "name_en", "name_tc", "name_sc", "bg_color", "text_color", "is_active", "created_at", "updated_at"}

func GetCompanyById(ctx context.Context, dbpool *pgx.Conn, id int) (*Company, error) {
	query := fmt.Sprintf("SELECT %s FROM tbl_company WHERE id = %d",
		strings.Join(companyColumns, ","), id)
	rows, err := dbpool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var company Company
	for rows.Next() {
		err := rows.Scan(
			&company.ID,
			&company.SlugEn,
			&company.SlugTc,
			&company.SlugSc,
			&company.NameEn,
			&company.NameTc,
			&company.NameSc,
			&company.BgColor,
			&company.TextColor,
			&company.IsActive,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return &company, nil
}
