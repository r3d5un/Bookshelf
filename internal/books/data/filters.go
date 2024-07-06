package data

import (
	"time"

	"github.com/google/uuid"
)

type MetaData struct {
	ItemsPerPage int `json:"itemsPerPage,omitempty"`
	TotalResults int `json:"totalResults,omitempty"`
	StartIndex   int `json:"startIndex,omitempty"`
}

type Filters struct {
	Page            int        `json:"page,omitempty"`
	PageSize        int        `json:"pageSize,omitempty"`
	StartIndex      int        `json:"startIndex,omitempty"`
	Count           int        `json:"count,omitempty"`
	ID              *uuid.UUID `json:"id,omitempty"`
	AuthorID        *uuid.UUID `json:"authorId,omitempty"`
	Title           *string    `json:"title,omitempty"`
	Description     *string    `json:"description,omitempty"`
	Name            *string    `json:"name,omitempty"`
	PublishedFrom   *time.Time `json:"publishedFrom,omitempty"`
	PublishedTo     *time.Time `json:"publishedTo,omitempty"`
	CreatedAtFrom   *time.Time `json:"createdAtFrom,omitempty"`
	CreatedAtTo     *time.Time `json:"createdAtTo,omitempty"`
	UpdatedAtFrom   *time.Time `json:"updatedAtFrom,omitempty"`
	UpdatedAtTo     *time.Time `json:"updatedAtTo,omitempty"`
	OrderBy         []string   `json:"order_by,omitempty"`
	OrderBySafeList []string   `json:"order_by_safe_list,omitempty"`
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
