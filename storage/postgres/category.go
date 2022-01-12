package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
)

type categoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CategoryCreate(category pb.Category) (pb.Category, error) {
	parId := stringToNullString(category.ParentId)

	err := r.db.QueryRow(`
		INSERT INTO categories (id, name, parent_id) 
		VALUES($1, $2, $3) returning id`,
		category.Id,
		category.Name,
		parId,
	).Scan(&category.Id)
	if err != nil {
		return pb.Category{}, err
	}
	category, err = r.CategoryGet(category.Id)
	if err != nil {
		return pb.Category{}, err
	}

	return category, nil
}

func (r *categoryRepo) CategoryGet(id string) (pb.Category, error) {
	var (
		parId    sql.NullString
		category pb.Category
	)

	err := r.db.QueryRow(`
		SELECT id, name, parent_id FROM categories WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&category.Id,
		&category.Name,
		&parId)
	if err != nil {
		return pb.Category{}, err
	}

	if parId.Valid {
		category.ParentId = parId.String
	}

	return category, nil
}

func (r *categoryRepo) CategoryList(page, limit int64) ([]*pb.Category, int64, error) {
	offset := (page - 1) * limit
	var parId sql.NullString
	rows, err := r.db.Queryx(`
		SELECT id, name, parent_id FROM categories WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		categories []*pb.Category
		count      int64
	)

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&parId)
		if parId.Valid {
			category.ParentId = parId.String
		}

		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, &category)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM categories WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return categories, count, nil
}

func (r *categoryRepo) CategoryUpdate(category pb.Category) (pb.Category, error) {
	parID := stringToNullString(category.ParentId)

	result, err := r.db.Exec(`UPDATE categories SET name=$1, parent_id=$2, updated_at=$3 WHERE id=$4`,
		&category.Name,
		&parID,
		time.Now().UTC(),
		&category.Id)
	if err != nil {
		return pb.Category{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Category{}, sql.ErrNoRows
	}

	category, err = r.CategoryGet(category.Id)
	if err != nil {
		return pb.Category{}, err
	}

	return category, err
}

func (r *categoryRepo) CategoryDelete(id string) error {
	result, err := r.db.Exec(`UPDATE categories SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func stringToNullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.Valid = true
		ns.String = s
		return ns
	}

	return ns
}
