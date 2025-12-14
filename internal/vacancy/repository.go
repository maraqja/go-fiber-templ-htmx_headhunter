package vacancy

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type RepositoryDI struct {
	DB *pgxpool.Pool
}

type PostgresRepository struct {
	db     *pgxpool.Pool
	logger *zerolog.Logger
}

func NewPostgresRepository(di RepositoryDI) *PostgresRepository {
	logger := log.Logger.With().Str("component", "VacancyPostgresRepository").Logger()
	return &PostgresRepository{
		db:     di.DB,
		logger: &logger,
	}
}

func (r *PostgresRepository) CreateVacancy(form VacancyCreateForm) {

}
