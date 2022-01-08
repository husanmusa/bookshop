package postgres

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/husanmus/bookshop/genproto/catalogService"
)

type catalogRepo struct {
	db *sqlx.DB
}

func NewCatalogRepo(db *sqlx.DB) *catalogRepo {
	return &catalogRepo{db: db}
}

func (r *catalogRepo) AuthorCreate(author pb.Author) (pb.Author, error) {
	var id string

	err := r.db.QueryRow(`
		INSERT INTO authors(id, name) 
		VALUES($1, $2) returning id`,
		author.Id,
		author.Name,
	).Scan(&id)

	if err != nil {
		return pb.Author{}, err
	}

	author, err = r.AuthorGet(id)
	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *catalogRepo) AuthorGet(id string) (pb.Author, error) {
	var author pb.Author
	err := r.db.QueryRow(`
		SELECT id, name FROM authors WHERE id = $1 AND delete_at IS NULL`, id).Scan(
		&author.Id,
		&author.Name)
	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *catalogRepo) AuthorList(page, limit int64) ([]*pb.Author, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`
		SELECT id, name FROM authors WHERE delete_at IS NULL LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		authors []*pb.Author
		count   int64
	)

	for rows.Next() {
		var author pb.Author

		err = rows.Scan(
			&author.Id,
			&author.Name)
		if err != nil {
			return nil, 0, err
		}
		authors = append(authors, &author)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM authors WHERE delete_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return authors, count, nil
}

func (r *catalogRepo) AuthorUpdate(author pb.Author) (pb.Author, error) {
	result, err := r.db.Exec(`UPDATE author SET name=$1 update_at=current_timestamp WHERE id=$2`,
		&author.Name)
	if err != nil {
		return pb.Author{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Author{}, sql.ErrNoRows
	}

	author, err = r.AuthorGet(author.Id)
	if err != nil {
		return pb.Author{}, err
	}

	return author, err
}

func (r *catalogRepo) AuthorDelete(id string) error {
	result, err := r.db.Exec(`UPDATE auhtors SET deleted_at = $1 WHERE id = $2`, time.Now(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
