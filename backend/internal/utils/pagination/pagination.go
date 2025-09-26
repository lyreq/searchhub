package pagination

import (
	"fmt"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Strategy int

const (
	None Strategy = iota
	Offset
	Cursor
)

type Pagination struct {
	Strategy

	// Offset pagination
	Page int

	// Cursor pagination
	Cursor uint

	// Shared
	Direction int // forward: 1, reverse: -1 (default: 1)
	Count     int // default: 10
}

func (data *Pagination) String() string {
	if data.Strategy == Offset {
		return fmt.Sprintf("%v-%v-%v", Offset, data.Page, data.Count)
	}
	if data.Strategy == Cursor {
		return fmt.Sprintf("%v-%v-%v-%v", Cursor, data.Cursor, data.Direction, data.Count)
	}

	return ""
}

func GetData(c echo.Context) Pagination {
	data := Pagination{
		Strategy:  None,
		Page:      1,
		Cursor:    0,
		Direction: 1,
		Count:     10,
	}

	count, err := strconv.Atoi(c.QueryParam("count"))
	if err == nil {
		data.Count = count
	}

	cursor := c.QueryParam("cursor")
	if cursor != "" {
		data.Strategy = Cursor
		cursor, err := strconv.Atoi(c.QueryParam("cursor"))
		if err == nil {
			data.Cursor = uint(cursor)
		}
		direction, err := strconv.Atoi(c.QueryParam("direction"))
		if err == nil {
			data.Direction = direction
		}

		return data
	}

	data.Strategy = Offset
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err == nil {
		data.Page = int(math.Max(1, float64(page)))
	}

	return data
}

func Prepare(db *gorm.DB, paginate *Pagination) *gorm.DB {
	if paginate == nil {
		return db
	}

	if paginate.Strategy == Cursor {
		db = db.Where("id > ?", paginate.Cursor)
	} else if paginate.Strategy == Offset {
		db = db.Offset((paginate.Page - 1) * paginate.Count)
	}

	if paginate.Direction == 1 {
		db = db.Order("id ASC")
	} else if paginate.Direction == -1 {
		db = db.Order("id DESC")
	}

	db = db.Limit(paginate.Count)
	return db
}

func (p *Pagination) Paginate(db *gorm.DB) *gorm.DB {
	return Prepare(db, p)
}
