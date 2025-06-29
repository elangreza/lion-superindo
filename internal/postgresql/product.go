package postgresql

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/elangreza/lion-superindo/internal/domain"
	"github.com/elangreza/lion-superindo/internal/params"
)

func (pr *PostgresRepo) listQuery(req params.ListProductsQueryParams) squirrel.SelectBuilder {
	q := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).Select().From("products p")

	search := strings.TrimSpace(req.Search)
	if len(search) != 0 {
		if _, err := strconv.Atoi(search); err == nil {
			q = q.Where(squirrel.Eq{"p.id": search})
		} else {
			q = q.Where(squirrel.Like{"LOWER(p.name)": "%" + strings.ToLower(search) + "%"})
		}
	}

	if len(req.Types) != 0 {
		q = q.Where(squirrel.Eq{"p.product_type_name": req.Types})
	}

	return q
}

func (pr *PostgresRepo) ListProducts(ctx context.Context, req params.ListProductsQueryParams) ([]domain.Product, error) {
	q := pr.listQuery(req).Columns("id", "name", "price", "product_type_name", "created_at")

	if req.GetSortMapping() != nil {
		for key, direction := range req.GetSortMapping() {
			q = q.OrderBy(key + " " + direction)
		}
	} else {
		q = q.OrderBy("id asc")
	}

	q = q.Limit(uint64(req.Limit))
	if req.Page > 1 {
		q = q.Offset(uint64(req.Limit * (req.Page - 1)))
	}

	qr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pr.db.QueryContext(ctx, qr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.ProductType.Name,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *PostgresRepo) CountProducts(ctx context.Context, req params.ListProductsQueryParams) (int, error) {
	qCount := pr.listQuery(req).Columns("count(id)")
	qc, args, err := qCount.ToSql()
	if err != nil {
		return 0, err
	}

	var countProducts int
	if err = pr.db.QueryRow(qc, args...).Scan(&countProducts); err != nil {
		return 0, err
	}

	return countProducts, nil
}

func (pr *PostgresRepo) CreateProduct(ctx context.Context, req params.CreateProductRequest) (int, error) {
	var id int
	err := runInTx(ctx, pr.db, func(tx *sql.Tx) error {
		qInsertProductType := `INSERT INTO product_types("name") VALUES($1) ON CONFLICT(name) DO NOTHING;`
		if _, err := tx.ExecContext(ctx, qInsertProductType, req.Type); err != nil {
			return err
		}

		qInsertProduct :=
			`INSERT INTO products("name", price, product_type_name) VALUES($1, $2, $3) RETURNING id;`
		if err := tx.QueryRowContext(ctx, qInsertProduct, req.Name, req.Price, req.Type).Scan(&id); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}
