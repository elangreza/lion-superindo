package postgresql

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/elangreza14/superindo/internal/domain"
	"github.com/elangreza14/superindo/internal/params"
	"github.com/redis/go-redis/v9"
)

type (
	ProductRepo struct {
		db    *sql.DB
		cache *redis.Client
	}
)

func NewProductRepo(db *sql.DB, cache *redis.Client) *ProductRepo {
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

func (pr *ProductRepo) ListProduct(ctx context.Context, req params.ListProductQueryParams) (products []domain.Product, err error) {
	keyRaw := req.GetKey()
	key := "listProduct:" + string(keyRaw)

	rcmd := pr.cache.Get(ctx, key)
	if err != nil && err != redis.Nil {
		return
	}
	if rcmd.Val() != "" {
		err = json.Unmarshal([]byte(rcmd.Val()), &products)
		if err != nil {
			return
		}
		slog.Info("using redis", "method", "ListProduct")
		return
	}

	q := pr.ListQuery(req).Columns("id", "name", "price", "product_type_name", "created_at", "updated_at")

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

	// products := []domain.Product{}
	for rows.Next() {
		product := domain.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
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

	byteProducts, err := json.Marshal(products)
	if err != nil {
		return
	}
	err = pr.cache.Set(ctx, key, string(byteProducts), time.Second*60).Err()
	if err != nil {
		return
	}

	return products, nil
}

func (pr *ProductRepo) TotalProduct(ctx context.Context, req params.ListProductQueryParams, withCache bool) (totalProducts int, err error) {
	keyRaw := req.GetKey()
	key := "totalProduct:" + string(keyRaw)

	if withCache {
		rcmd := pr.cache.Get(ctx, key)
		if rcmd.Val() != "" {
			err = rcmd.Scan(&totalProducts)
			if err != nil && err != redis.Nil {
				return
			}

			slog.Info("using redis", "method", "TotalProduct")
			return
		}
	}

	qCount := pr.ListQuery(req).Columns("count(id)")
	qc, args, err := qCount.ToSql()
	if err != nil {
		return
	}

	if err = pr.db.QueryRow(qc, args...).Scan(&totalProducts); err != nil {
		return
	}

	if withCache {
		if err = pr.cache.Set(ctx, key, totalProducts, time.Second*60).Err(); err != nil {
			return
		}
	}

	return
}

func (pr *ProductRepo) CreateProduct(ctx context.Context, req params.CreateProductRequest) (id int, err error) {

	// id := 0
	err = runInTx(ctx, pr.db, func(tx *sql.Tx) error {
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
		return
	}

	if err = pr.cache.FlushAll(ctx).Err(); err != nil {
		return
	}

	return
}
