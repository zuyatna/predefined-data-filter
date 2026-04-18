package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"predefined-data-filter/internal/domain"
)

type ProductHandler struct {
	useCase domain.ProductUseCase
}

func NewProductHandler(mux *http.ServeMux, useCase domain.ProductUseCase) {
	handler := &ProductHandler{
		useCase: useCase,
	}

	// Using standard http.ServeMux (Go 1.22 supports method routing)
	mux.HandleFunc("GET /api/v1/products", handler.FetchProducts)
}

func (h *ProductHandler) FetchProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	filter := buildFilterFromQuery(r)

	res, err := h.useCase.FetchProducts(r.Context(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func buildFilterFromQuery(r *http.Request) domain.ProductFilter {
	var filter domain.ProductFilter

	query := r.URL.Query()

	if q := query.Get("search_query"); q != "" {
		filter.SearchQuery = q
	}

	if c := query.Get("category_id"); c != "" {
		if id, err := strconv.Atoi(c); err == nil {
			filter.CategoryID = &id
		}
	}

	if min := query.Get("min_price"); min != "" {
		if val, err := strconv.ParseFloat(min, 64); err == nil {
			filter.MinPrice = &val
		}
	}

	if max := query.Get("max_price"); max != "" {
		if val, err := strconv.ParseFloat(max, 64); err == nil {
			filter.MaxPrice = &val
		}
	}

	if colors := query.Get("color_id"); colors != "" {
		parts := strings.Split(colors, ",")
		for _, p := range parts {
			if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				filter.ColorIDs = append(filter.ColorIDs, id)
			}
		}
	}

	if labels := query.Get("label_id"); labels != "" {
		parts := strings.Split(labels, ",")
		for _, p := range parts {
			if id, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
				filter.LabelIDs = append(filter.LabelIDs, id)
			}
		}
	}

	if start := query.Get("start_date"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			filter.StartDate = &t
		}
	}

	if end := query.Get("end_date"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			// Add 23:59:59 to include the whole end day
			t = t.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
			filter.EndDate = &t
		}
	}

	if sort := query.Get("sort"); sort != "" {
		filter.Sort = sort
	}

	if page := query.Get("page"); page != "" {
		if val, err := strconv.Atoi(page); err == nil {
			filter.Pagination.Page = val
		}
	}

	if limit := query.Get("limit"); limit != "" {
		if val, err := strconv.Atoi(limit); err == nil {
			filter.Pagination.Limit = val
		}
	}

	return filter
}
