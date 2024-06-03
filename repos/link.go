package repos

import (
	"database/sql"
	"time"

	"github.com/aacebo/sss/models"
)

type ILinkRepository interface {
	GetOne(fromId string, toId string) (models.Link, bool)
	GetByPageID(pageId string) []models.Link

	Create(value models.Link) models.Link
	DeleteByPageID(pageId string)
}

type LinkRepository struct {
	pg *sql.DB
}

func NewLink(pg *sql.DB) LinkRepository {
	return LinkRepository{pg}
}

func (self LinkRepository) GetOne(fromId string, toId string) (models.Link, bool) {
	v := models.Link{}
	err := self.pg.QueryRow(
		`
			SELECT
				from_id,
				to_id,
				created_at
			FROM links
			WHERE from_id = $1
			AND to_id = $2
		`,
		fromId,
		toId,
	).Scan(
		&v.FromID,
		&v.ToID,
		&v.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return v, false
		}

		panic(err)
	}

	return v, true
}

func (self LinkRepository) GetByPageID(pageId string) []models.Link {
	rows, err := self.pg.Query(
		`
			SELECT
				from_id,
				to_id,
				created_at
			FROM links
			WHERE from_id = $1
		`,
		pageId,
	)

	if err != nil {
		panic(err)
	}

	defer rows.Close()
	arr := []models.Link{}

	for rows.Next() {
		v := models.Link{}
		err := rows.Scan(
			&v.FromID,
			&v.ToID,
			&v.CreatedAt,
		)

		if err != nil {
			panic(err)
		}

		arr = append(arr, v)
	}

	return arr
}

func (self LinkRepository) Create(value models.Link) models.Link {
	now := time.Now()
	value.CreatedAt = now
	_, err := self.pg.Exec(
		`
			INSERT INTO links (
				from_id,
				to_id,
				created_at
			) VALUES (
				$1,
				$2,
				$3
			)
		`,
		value.FromID,
		value.ToID,
		value.CreatedAt,
	)

	if err != nil {
		panic(err)
	}

	return value
}

func (self LinkRepository) DeleteByPageID(pageId string) {
	_, err := self.pg.Exec(
		`
			DELETE FROM links
			WHERE from_id = $1
		`,
		pageId,
	)

	if err != nil {
		panic(err)
	}
}
