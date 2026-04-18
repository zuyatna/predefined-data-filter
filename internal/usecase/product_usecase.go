package usecase

import (
	"context"
	"math"
	"predefined-data-filter/internal/domain"
)

type productUseCase struct {
	repo domain.ProductRepository
}

func NewProductUseCase(repo domain.ProductRepository) domain.ProductUseCase {
	return &productUseCase{repo: repo}
}

func (u *productUseCase) FetchProducts(ctx context.Context, filter domain.ProductFilter) (domain.PaginatedProductResponse, error) {
	// 1. Validate and set default pagination
	if filter.Pagination.Page <= 0 {
		filter.Pagination.Page = 1
	}
	if filter.Pagination.Limit <= 0 {
		filter.Pagination.Limit = 10
	}
	if filter.Pagination.Limit > 100 { // Max limit to prevent overload
		filter.Pagination.Limit = 100
	}

	// 2. Fetch from repository
	products, totalItems, err := u.repo.Fetch(ctx, filter)
	if err != nil {
		return domain.PaginatedProductResponse{}, err
	}

	// 3. Calculate pagination metadata
	totalPages := int(math.Ceil(float64(totalItems) / float64(filter.Pagination.Limit)))

	return domain.PaginatedProductResponse{
		Data: products,
		Pagination: domain.PaginationResponse{
			CurrentPage:  filter.Pagination.Page,
			ItemsPerPage: filter.Pagination.Limit,
			TotalItems:   totalItems,
			TotalPages:   totalPages,
		},
	}, nil
}
