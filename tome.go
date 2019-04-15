// Copyright (c) 2019 Cyro Dubeux. License MIT.

// Package tome was designed to paginate simple RESTful APIs.
package tome

import (
	"errors"
	"fmt"
	"math"
)

// Chapter type is a struct for pagination results.
type Chapter struct {
	// Data that you want to return along with pagination settings.
	Data interface{} `json:"data"`
	// API base URL.
	BaseURL string `json:"base_url"`
	// The first URL link with page number.
	FirstURL string `json:"first_url"`
	// The next URL link with page number.
	NextURL string `json:"next_url"`
	// The previous URL link with page number.
	PreviousURL string `json:"prev_url"`
	// The last URL link with page number.
	LastURL string `json:"last_url"`
	// The inicial offset position.
	Offset int `json:"-"`
	// Limit per page.
	Limit int `json:"per_page"`
	// The page number captured on the request params.
	NewPage int `json:"-"`
	// Current page of the tome.
	CurrentPage int `json:"current_page"`
	// The last page of the tome.
	LastPage int `json:"last_page"`
	// Total of pages, this usually comes from a SQL query total rows result.
	TotalPages int `json:"total"`
}

// Paginate handles the pagination calculation.
func (c *Chapter) Paginate() (*Chapter, error) {
	if c.BaseURL == "" {
		return nil, errors.New("Base URL is missing")
	}

	c.setDefaults()                 // Checking if need defaults
	c.ceilLastPage()                // Ceiling the last page.
	offset, limit := c.doPaginate() // Pagination calculation.
	c.createLinks()                 // Creating links.

	return &Chapter{
		c.Data,
		c.BaseURL,
		c.FirstURL,
		c.NextURL,
		c.PreviousURL,
		c.LastURL,
		offset,
		limit,
		c.NewPage,
		c.CurrentPage,
		c.LastPage,
		c.TotalPages,
	}, nil
}

// Calculates the offset and the limit.
func (c *Chapter) doPaginate() (int, int) {
	if c.NewPage > c.CurrentPage {
		c.CurrentPage = c.NewPage
		c.Offset = (c.CurrentPage - 1) * c.Limit
	}
	return c.Offset, c.Limit
}

// Ceils the last page and generates
// a integer number.
func (c *Chapter) ceilLastPage() {
	c.LastPage = int(math.Ceil(float64(c.TotalPages) / float64(c.Limit)))
}

// Creates next and previous links using
// the given base URL.
func (c *Chapter) createLinks() {
	c.FirstURL = fmt.Sprintf("%s?page=%d", c.BaseURL, 1)
	if c.NewPage < c.LastPage {
		c.NextURL = fmt.Sprintf("%s?page=%d", c.BaseURL, c.CurrentPage+1)
	}
	if c.LastPage > c.NewPage {
		c.PreviousURL = fmt.Sprintf("%s?page=%d", c.BaseURL, c.CurrentPage-1)
	}
	c.LastURL = fmt.Sprintf("%s?page=%d", c.BaseURL, c.LastPage)
}

// Sets the defaults values for current page
// and limit if none of them were provided.
func (c *Chapter) setDefaults() {
	if cp := c.CurrentPage == 0; cp {
		c.CurrentPage = 1
	}
	if l := c.Limit == 0; l {
		c.Limit = 10
	}
}
