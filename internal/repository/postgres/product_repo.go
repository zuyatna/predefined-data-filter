package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"predefined-data-filter/internal/domain"
	"github.com/lib/pq"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Fetch(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	// 1. Build Base Query
	query := `
		SELECT p.id, p.name, p.price, p.purchases_count, p.reviews_count, p.created_at,
		       c.id, c.name
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE 1=1
	`
	countQuery := `SELECT COUNT(p.id) FROM products p WHERE 1=1`

	args := []interface{}{}
	argId := 1

	var conditions []string

	// 2. Apply Filters dynamically
	if filter.SearchQuery != "" {
		conditions = append(conditions, fmt.Sprintf("p.name ILIKE $%d", argId))
		args = append(args, "%"+filter.SearchQuery+"%")
		argId++
	}

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("p.category_id = $%d", argId))
		args = append(args, *filter.CategoryID)
		argId++
	}

	if filter.MinPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price >= $%d", argId))
		args = append(args, *filter.MinPrice)
		argId++
	}

	if filter.MaxPrice != nil {
		conditions = append(conditions, fmt.Sprintf("p.price <= $%d", argId))
		args = append(args, *filter.MaxPrice)
		argId++
	}

	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at >= $%d", argId))
		args = append(args, *filter.StartDate)
		argId++
	}

	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("p.created_at <= $%d", argId))
		args = append(args, *filter.EndDate)
		argId++
	}

	if len(filter.ColorIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM product_colors pc WHERE pc.product_id = p.id AND pc.color_id = ANY($%d))", argId))
		args = append(args, pq.Array(filter.ColorIDs))
		argId++
	}

	if len(filter.LabelIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM product_labels pl WHERE pl.product_id = p.id AND pl.label_id = ANY($%d))", argId))
		args = append(args, pq.Array(filter.LabelIDs))
		argId++
	}

	// Combine conditions
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		query += whereClause
		countQuery += whereClause
	}

	// 3. Execute Count Query First
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting products: %w", err)
	}

	if total == 0 {
		return []domain.Product{}, 0, nil
	}

	// 4. Apply Sorting
	switch filter.Sort {
	case "popular":
		query += " ORDER BY (p.purchases_count + p.reviews_count) DESC"
	case "newest":
		query += " ORDER BY p.created_at DESC"
	case "price_asc":
		query += " ORDER BY p.price ASC"
	case "price_desc":
		query += " ORDER BY p.price DESC"
	default:
		query += " ORDER BY p.id ASC"
	}

	// 5. Apply Pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argId, argId+1)
	args = append(args, filter.Pagination.Limit, filter.Pagination.Offset())

	// 6. Execute Main Query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	var productIDs []int
	productMap := make(map[int]*domain.Product)

	for rows.Next() {
		var p domain.Product
		var catID sql.NullInt64
		var catName sql.NullString

		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.PurchasesCount, &p.ReviewsCount, &p.CreatedAt,
			&catID, &catName,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning product row: %w", err)
		}

		if catID.Valid {
			p.Category.ID = int(catID.Int64)
			p.Category.Name = catName.String
		}
		
		p.Colors = []domain.Color{}
		p.Labels = []domain.Label{}

		products = append(products, p)
		productIDs = append(productIDs, p.ID)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating product rows: %w", err)
	}

	for i := range products {
		productMap[products[i].ID] = &products[i]
	}

	// 7. Fetch Relationships (Colors and Labels)
	if len(productIDs) > 0 {
		// Get Colors
		err = r.fetchColorsForProducts(ctx, productIDs, productMap)
		if err != nil {
			return nil, 0, fmt.Errorf("error fetching colors: %w", err)
		}

		// Get Labels
		err = r.fetchLabelsForProducts(ctx, productIDs, productMap)
		if err != nil {
			return nil, 0, fmt.Errorf("error fetching labels: %w", err)
		}
	}

	return products, total, nil
}

func (r *productRepository) fetchColorsForProducts(ctx context.Context, productIDs []int, productMap map[int]*domain.Product) error {
	query := `
		SELECT pc.product_id, c.id, c.name
		FROM product_colors pc
		JOIN colors c ON pc.color_id = c.id
		WHERE pc.product_id = ANY($1)
	`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(productIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int
		var color domain.Color
		if err := rows.Scan(&productID, &color.ID, &color.Name); err != nil {
			return err
		}
		if p, ok := productMap[productID]; ok {
			p.Colors = append(p.Colors, color)
		}
	}
	return rows.Err()
}

func (r *productRepository) fetchLabelsForProducts(ctx context.Context, productIDs []int, productMap map[int]*domain.Product) error {
	query := `
		SELECT pl.product_id, l.id, l.name
		FROM product_labels pl
		JOIN labels l ON pl.label_id = l.id
		WHERE pl.product_id = ANY($1)
	`
	rows, err := r.db.QueryContext(ctx, query, pq.Array(productIDs))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID int
		var label domain.Label
		if err := rows.Scan(&productID, &label.ID, &label.Name); err != nil {
			return err
		}
		if p, ok := productMap[productID]; ok {
			p.Labels = append(p.Labels, label)
		}
	}
	return rows.Err()
}
