package postgresql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type ProductRepo struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{db}
}

func (pr *ProductRepo) ListQuery(req params.ProductQueryParams) squirrel.SelectBuilder {
	q := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select().From("products p")

	if len(req.Search) != 0 {
		if _, err := strconv.Atoi(req.Search); err == nil {
			q = q.Where(squirrel.Eq{"p.id": req.Search})
		} else {
			q = q.Where(squirrel.Like{"LOWER(p.name)": "%" + strings.ToLower(req.Search) + "%"})
		}
	}

	if len(req.Types) != 0 {
		q = q.Where(squirrel.Eq{"p.product_type_name": req.Types})
	}

	return q
}

func (pr *ProductRepo) ListProduct(ctx context.Context, req params.ProductQueryParams) (int, []domain.Product, error) {
	qBase := pr.ListQuery(req)
	qCount := qBase.Columns("count(id)")
	qc, args, err := qCount.ToSql()
	if err != nil {
		return 0, nil, err
	}
	totalProducts := 0
	err = pr.db.QueryRow(qc, args...).Scan(&totalProducts)
	if err != nil {
		return 0, nil, err
	}

	if totalProducts == 0 {
		return 0, nil, nil
	}

	q := qBase.Columns("id", "name", "quantity", "price", "product_type_name", "created_at", "updated_at")

	if req.GetSortMapping() != nil {
		for key, direction := range req.GetSortMapping() {
			if key == "updated_at" {
				q = q.OrderBy("coalesce(p.updated_at, p.created_at)" + " " + direction)
			} else {
				q = q.OrderBy(key + " " + direction)
			}
		}
	} else {
		q = q.OrderBy("id asc")
	}

	q = q.Limit(uint64(req.Limit))
	if req.Page > 1 {
		q = q.Offset(req.Limit * (req.Page - 1))
	}

	qr, args, err := q.ToSql()
	if err != nil {
		return 0, nil, err
	}

	rows, err := pr.db.QueryContext(ctx, qr, args...)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Quantity,
			&product.Price,
			&product.ProductType.Name,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return 0, nil, err
		}
		products = append(products, product)
	}

	return totalProducts, products, nil
}
