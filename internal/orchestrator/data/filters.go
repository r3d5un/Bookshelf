package data

import (
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Metadata struct {
	CurrentPage  int    `json:"current_page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	FirstPage    int    `json:"first_page,omitempty"`
	LastPage     int    `json:"last_page,omitempty"`
	TotalRecords int    `json:"total_records,omitempty"`
	OrderBy      string `json:"order_by,omitempty"`
}

type Filters struct {
	Page            int        `json:"page,omitempty"`
	PageSize        int        `json:"pageSize,omitempty"`
	StartIndex      int        `json:"startIndex,omitempty"`
	Count           int        `json:"count,omitempty"`
	ID              *uuid.UUID `json:"id,omitempty"`
	Queue           *string    `json:"queue"`
	State           *string    `json:"state,omitempty"`
	CreatedAtFrom   *time.Time `json:"createdAtFrom,omitempty"`
	CreatedAtTo     *time.Time `json:"createdAtTo,omitempty"`
	UpdatedAtFrom   *time.Time `json:"updatedAtFrom,omitempty"`
	UpdatedAtTo     *time.Time `json:"updatedAtTo,omitempty"`
	RunAtFrom       *time.Time `json:"runAtFrom,omitempty"`
	RunAtTo         *time.Time `json:"runAtTo,omitempty"`
	OrderBy         []string   `json:"orderBy,omitempty"`
	OrderBySafeList []string   `json:"orderBySafeList,omitempty"`
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

func calculateMetadata(totalRecords, page, pageSize int, orderBySlice []string) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
		OrderBy:      strings.Join(orderBySlice, ","),
	}
}
