package repos

import (
	"database/sql"
	"time"

	"github.com/aacebo/sss/models"
)

type IPageRepository interface {
	GetByID(id string) (models.Page, bool)
	GetOne(url string) (models.Page, bool)

	Create(value models.Page) models.Page
	Update(value models.Page) models.Page
}

type PageRepository struct {
	pg *sql.DB
}

func NewPage(pg *sql.DB) PageRepository {
	return PageRepository{pg}
}

func (self PageRepository) GetByID(id string) (models.Page, bool) {
	v := models.Page{}
	err := self.pg.QueryRow(
		`
			SELECT
				id,
				domain_id,
				title,
				url,
				address,
				size,
				elapse_ms,
				link_count,
				created_at,
				updated_at
			FROM pages
			WHERE id = $1
		`,
		id,
	).Scan(
		&v.ID,
		&v.DomainID,
		&v.Title,
		&v.Url,
		&v.Address,
		&v.Size,
		&v.ElapseMs,
		&v.LinkCount,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return v, false
		}

		panic(err)
	}

	return v, true
}

func (self PageRepository) GetOne(url string) (models.Page, bool) {
	v := models.Page{}
	err := self.pg.QueryRow(
		`
			SELECT
				id,
				domain_id,
				title,
				url,
				address,
				size,
				elapse_ms,
				link_count,
				created_at,
				updated_at
			FROM pages
			WHERE url = $1
		`,
		url,
	).Scan(
		&v.ID,
		&v.DomainID,
		&v.Title,
		&v.Url,
		&v.Address,
		&v.Size,
		&v.ElapseMs,
		&v.LinkCount,
		&v.CreatedAt,
		&v.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return v, false
		}

		panic(err)
	}

	return v, true
}

func (self PageRepository) Create(value models.Page) models.Page {
	now := time.Now()
	value.CreatedAt = now
	value.UpdatedAt = now
	_, err := self.pg.Exec(
		`
			INSERT INTO pages (
				id,
				domain_id,
				title,
				url,
				address,
				size,
				elapse_ms,
				link_count,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10
			)
		`,
		value.ID,
		value.DomainID,
		value.Title,
		value.Url,
		value.Address,
		value.Size,
		value.ElapseMs,
		value.LinkCount,
		value.CreatedAt,
		value.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	return value
}

func (self PageRepository) Update(value models.Page) models.Page {
	now := time.Now()
	value.UpdatedAt = now
	_, err := self.pg.Exec(
		`
			UPDATE pages SET
				title = $2,
				address = $3,
				size = $4,
				elapse_ms = $5,
				link_count = $6,
				updated_at = $7
			WHERE id = $1
		`,
		value.ID,
		value.Title,
		value.Address,
		value.Size,
		value.ElapseMs,
		value.LinkCount,
		value.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	return value
}
