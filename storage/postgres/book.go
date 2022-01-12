package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	pb "github.com/husanmusa/bookshop/genproto/catalog"
	"github.com/husanmusa/bookshop/pkg/utils"
)

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepo(db *sqlx.DB) *bookRepo {
	return &bookRepo{db: db}
}

func (r *bookRepo) BookCreate(book pb.Book) (pb.Book, error) {
	err := r.db.QueryRow(`
		INSERT INTO books(id, name, author_id) 
		VALUES($1, $2, $3) returning id`,
		book.Id,
		book.Name,
		book.AuthorId.Id,
	).Scan(&book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	for _, v := range book.CategoryId {
		_, err = r.db.Queryx(`
		INSERT INTO book_categories(book_id, category_id) 
		VALUES($1, $2)`,
			book.Id,
			v,
		)
		if err != nil {
			return pb.Book{}, err
		}
	}

	book, err = r.BookGet(book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	return book, nil
}

func (r *bookRepo) BookGet(id string) (pb.Book, error) {
	var book pb.Book
	var name pb.Author
	err := r.db.QueryRow(`
	select b.id, b.name, a.name, b.created_at, b.updated_at from books b
	JOIN authors a on b.author_id = a.id
    WHERE b.id=$1 AND b.deleted_at IS NULL`, id).
		Scan(
			&book.Id,
			&book.Name,
			&name.Name,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
	book.AuthorId = &name
	if err != nil {
		return pb.Book{}, err
	}

	rows, err := r.db.Queryx(`
	select c.id ,c.name, c.parent_id, c.created_at, c.updated_at from book_categories as bc
	join categories c on bc.category_id = c.id
	where bc.book_id=$1`, id)
	if err != nil {
		return pb.Book{}, err
	}
	if err = rows.Err(); err != nil {
		return pb.Book{}, err
	}

	var parId sql.NullString

	for rows.Next() {
		var category pb.Category
		err = rows.Scan(
			&category.Id,
			&category.Name,
			&parId,
			&category.CreatedAt,
			&category.UpdatedAt)
		if err != nil {
			return pb.Book{}, err
		}
		if parId.Valid {
			category.ParentId = parId.String
		}

		book.Categories = append(book.Categories, &category)
	}

	defer rows.Close()

	return book, nil
}

func (r *bookRepo) BookList(page, limit int64, filters map[string]string) ([]*pb.Book, int64, error) {
	offset := (page - 1) * limit

	sb := sqlbuilder.NewSelectBuilder()

	sb.Select("b.id")
	sb.From("book_categories bc")
	sb.JoinWithOption("LEFT", "books b", "b.id=bc.book_id")
	if value, ok := filters["authors"]; ok {
		args := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sb.JoinWithOption("LEFT", "authors a", "a.id=b.author_id")
		sb.Where(sb.In("a.id", args...))
	}
	if value, ok := filters["categories"]; ok {
		args := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sb.JoinWithOption("LEFT", "categories c", "bc.category_id = c.id")
		sb.Where(sb.In("c.id", args...))
	}
	sb.Limit(int(limit))
	sb.Offset(int(offset))
	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var (
		books []*pb.Book
		count int64
	)
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, 0, err
		}
		book, _ := r.BookGet(id)
		books = append(books, &book)
	}

	sbc := sqlbuilder.NewSelectBuilder()
	sbc.Select("count(*)")
	sbc.From("book_categories bc")
	sbc.JoinWithOption("LEFT", "books b", "b.id=bc.book_id")
	if value, ok := filters["authors"]; ok {
		argsA := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sbc.JoinWithOption("LEFT", "authors a", "a.id=b.author_id")
		sbc.Where(sbc.In("a.id", argsA...))
	}
	if value, ok := filters["categories"]; ok {
		argsC := utils.StringSliceToInterfaceSlice(utils.ParseFilter(value))
		sbc.JoinWithOption("LEFT", "categories c", "bc.category_id = c.id")
		sbc.Where(sbc.In("c.id", argsC...))
	}
	query, args = sbc.BuildWithFlavor(sqlbuilder.PostgreSQL)
	fmt.Println(query, args)
	err = r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		fmt.Println("QABRISTAsa")
		return nil, 0, err
	}

	return books, count, nil
}

func (r *bookRepo) BookUpdate(book pb.Book) (pb.Book, error) {
	result, err := r.db.Exec(`UPDATE books SET name=$1, author_id=$2, updated_at=$3 WHERE id=$4`,
		&book.Name,
		&book.AuthorId.Id,
		time.Now().UTC(),
		&book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Book{}, sql.ErrNoRows
	}

	_, err = r.db.Exec(`DELETE from book_categories WHERE id=$1`, book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	for _, v := range book.CategoryId {
		_, err = r.db.Queryx(`
		INSERT INTO book_categories(book_id, category_id) 
		VALUES($1, $2)`,
			book.Id,
			v,
		)
		if err != nil {
			return pb.Book{}, err
		}
	}

	book, err = r.BookGet(book.Id)
	if err != nil {
		return pb.Book{}, err
	}

	return book, err
}

func (r *bookRepo) BookDelete(id string) error {
	result, err := r.db.Exec(`UPDATE books SET deleted_at = $1 WHERE id = $2`, time.Now().UTC(), id)
	if err != nil {
		return err
	}

	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
