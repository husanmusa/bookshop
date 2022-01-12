package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
)

type authorRepo struct {
	db *sqlx.DB
}

func NewAuthorRepo(db *sqlx.DB) *authorRepo {
	return &authorRepo{db: db}
}

func (r *authorRepo) AuthorCreate(author pb.Author) (pb.Author, error) {
	err := r.db.QueryRow(`
		INSERT INTO authors(id, name) 
		VALUES($1, $2) returning id`,
		author.Id,
		author.Name,
	).Scan(&author.Id)
	if err != nil {
		return pb.Author{}, err
	}
	author, err = r.AuthorGet(author.Id)
	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *authorRepo) AuthorGet(id string) (pb.Author, error) {
	fmt.Println(id)
	var author pb.Author
	err := r.db.QueryRow(`
		SELECT id, name FROM authors WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&author.Id,
		&author.Name)
	if err != nil {
		return pb.Author{}, err
	}

	return author, nil
}

func (r *authorRepo) AuthorList(page, limit int64) ([]*pb.Author, int64, error) {
	offset := (page - 1) * limit
	rows, err := r.db.Queryx(`
		SELECT id, name FROM authors WHERE deleted_at IS NULL LIMIT $1 OFFSET $2`, limit, offset)
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

	err = r.db.QueryRow(`SELECT count(*) FROM authors WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return authors, count, nil
}

func (r *authorRepo) AuthorUpdate(author pb.Author) (pb.Author, error) {
	result, err := r.db.Exec(`UPDATE authors SET name=$1, updated_at=$2 WHERE id=$3`,
		&author.Name,
		time.Now().UTC(),
		&author.Id)
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

func (r *authorRepo) AuthorDelete(id string) error {
	result, err := r.db.Exec(`UPDATE authors SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
