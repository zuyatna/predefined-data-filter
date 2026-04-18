package domain

import (
	"context"
	"time"
)

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Color struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Category       Category  `json:"category"`
	Price          float64   `json:"price"`
	PurchasesCount int       `json:"purchases_count"`
	ReviewsCount   int       `json:"reviews_count"`
	CreatedAt      time.Time `json:"created_at"`
	Colors         []Color   `json:"colors"`
	Labels         []Label   `json:"labels"`
}

// ProductFilter represents all available filter criteria
type ProductFilter struct {
	CategoryID   *int       `json:"category_id"`
	MinPrice     *float64   `json:"min_price"`
	MaxPrice     *float64   `json:"max_price"`
	ColorIDs     []int      `json:"color_ids"`
	LabelIDs     []int      `json:"label_ids"`
	StartDate    *time.Time `json:"start_date"`
	EndDate      *time.Time `json:"end_date"`
	Sort         string     `json:"sort"` // e.g., "popular", "newest", "price_asc", "price_desc"
	SearchQuery  string     `json:"search_query"`
	Pagination   PaginationRequest
}

type PaginationRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (p PaginationRequest) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
}

type PaginatedProductResponse struct {
	Data       []Product          `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

type ProductRepository interface {
	Fetch(ctx context.Context, filter ProductFilter) ([]Product, int, error)
}

type ProductUseCase interface {
	FetchProducts(ctx context.Context, filter ProductFilter) (PaginatedProductResponse, error)
}
