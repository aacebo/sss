package repos

import (
	"database/sql"
	"time"

	"github.com/aacebo/sss/models"
)

type IDomainRepository interface {
	GetByID(id string) (models.Domain, bool)
	GetOne(name string, extension string) (models.Domain, bool)

	Create(value models.Domain) models.Domain
	Update(value models.Domain) models.Domain
}

type DomainRepository struct {
	pg *sql.DB
}

func NewDomain(pg *sql.DB) DomainRepository {
	return DomainRepository{pg}
}

func (self DomainRepository) GetByID(id string) (models.Domain, bool) {
	v := models.Domain{}
	err := self.pg.QueryRow(
		`
			SELECT
				id,
				name,
				extension,
				created_at,
				updated_at
			FROM domains
			WHERE id = $1
		`,
		id,
	).Scan(
		&v.ID,
		&v.Name,
		&v.Extension,
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

func (self DomainRepository) GetOne(name string, extension string) (models.Domain, bool) {
	v := models.Domain{}
	err := self.pg.QueryRow(
		`
			SELECT
				id,
				name,
				extension,
				created_at,
				updated_at
			FROM domains
			WHERE name = $1
			AND extension = $2
		`,
		name,
		extension,
	).Scan(
		&v.ID,
		&v.Name,
		&v.Extension,
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

func (self DomainRepository) Create(value models.Domain) models.Domain {
	now := time.Now()
	value.CreatedAt = now
	value.UpdatedAt = now
	_, err := self.pg.Exec(
		`
			INSERT INTO domains (
				id,
				name,
				extension,
				created_at,
				updated_at
			) VALUES (
				$1,
				$2,
				$3,
				$4,
				$5
			)
		`,
		value.ID,
		value.Name,
		value.Extension,
		value.CreatedAt,
		value.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	return value
}

func (self DomainRepository) Update(value models.Domain) models.Domain {
	now := time.Now()
	value.UpdatedAt = now
	_, err := self.pg.Exec(
		`
			UPDATE domains SET
				updated_at = $2
			WHERE id = $1
		`,
		value.ID,
		value.UpdatedAt,
	)

	if err != nil {
		panic(err)
	}

	return value
}
