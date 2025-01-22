package postgresql

//go:generate mockgen -source $GOFILE -destination ../../mock/postgresql/mock_$GOFILE -package mock$GOPACKAGE

import (
	"context"
	"database/sql"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
)

type (
	Cache interface {
		Set(key string, Value any) error
	}

	ProductRepo struct {
		db    *sql.DB
		cache Cache
	}
)

func NewProductRepo(db *sql.DB, cache Cache) *ProductRepo {
	return &ProductRepo{db, cache}
}

func (pr *ProductRepo) ListQuery(req params.ListProductQueryParams) squirrel.SelectBuilder {
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

func (pr *ProductRepo) ListProduct(ctx context.Context, req params.ListProductQueryParams) ([]domain.Product, error) {

	q := pr.ListQuery(req).Columns("id", "name", "quantity", "price", "product_type_name", "created_at", "updated_at")

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
		return nil, err
	}

	rows, err := pr.db.QueryContext(ctx, qr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (pr *ProductRepo) TotalProduct(ctx context.Context, req params.ListProductQueryParams) (totalProducts int, err error) {
	qCount := pr.ListQuery(req).Columns("count(id)")
	qc, args, err := qCount.ToSql()
	if err != nil {
		return
	}

	if err = pr.db.QueryRow(qc, args...).Scan(&totalProducts); err != nil {
		return
	}

	return
}

func (pr *ProductRepo) CreateProduct(ctx context.Context, req params.CreateProductRequest) error {
	runInTx(ctx, pr.db, func(tx *sql.Tx) error {
		qInsertProductType := `INSERT INTO product_types("name") VALUES($1) ON CONFLICT(name) DO NOTHING;`
		if _, err := tx.ExecContext(ctx, qInsertProductType, req.Type); err != nil {
			return err
		}

		qInsertProduct := `INSERT INTO products("name", quantity, price, product_type_name) VALUES($1, $2, $3, $4);`
		if _, err := tx.ExecContext(ctx, qInsertProduct, req.Name, req.Quantity, req.Price, req.Type); err != nil {
			return err
		}

		return nil
	})

	return nil
}
