package models

import (
	"errors"
	"math"
	"strings"
)

type Filter struct {
	Page     int
	PageSize int
	OrderBy  string
	Query    string
}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	NextPage     int
	PrevPage     int
	LastPage     int
	TotalRecords int
}

func (f *Filter) Validate() error {

	if f.Page <= 0 || f.Page >= 10_000_000 {
		return errors.New("Invalid page range")
	}

	if f.PageSize <= 0 || f.PageSize >= 10_000_000 {
		return errors.New("Invalid page size")
	}

	return nil

}

func (f *Filter) addOrdering(query string) string {
	switch {
	case f.OrderBy == "popular":
		return strings.Replace(query, "#orderby#", "ORDER BY votes desc, p.created_at desc", 1)
	default:
		return strings.Replace(query, "#orderby#", "ORDER BY p.created_at desc", 1)
	}
}

func (f *Filter) addWhere(q string) string {
	if len(f.Query) > 0 {
		return strings.Replace(q, "#where#", "WHERE LOWER(p.title) LIKE $1", 1)
	}

	return strings.Replace(q, "#where#", "", 1)
}

func (f *Filter) addLimitOffset(q string) string {
	if len(f.Query) > 0 {
		return strings.Replace(q, "#limit#", "LIMIT $2 OFFSET $3", 1)
	}

	return strings.Replace(q, "#limit#", "LIMIT $1 OFFSET $2", 1)
}

func (f *Filter) applyTemplate(q string) string {
	return f.addLimitOffset(f.addWhere(f.addOrdering(q)))
}

func (f *Filter) limit() int {
	return f.PageSize
}

func (f *Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}

func CalculateMetadata(totalRecord, page, pageSize int) Metadata {
	if totalRecord == 0 {
		return Metadata{}
	}

	meta := Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecord) / float64(pageSize))),
		TotalRecords: totalRecord,
	}
	meta.NextPage = meta.CurrentPage + 1
	meta.PrevPage = meta.CurrentPage - 1

	if meta.CurrentPage <= meta.FirstPage {
		meta.PrevPage = 0

	}

	return meta
}
