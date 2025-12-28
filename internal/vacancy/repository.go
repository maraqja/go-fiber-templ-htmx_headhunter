package vacancy

import (
	"context"

	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) CreateVacancy(ctx context.Context, form VacancyCreateForm) error {
	query := `INSERT INTO vacancies (email, role, company, salary, type, location) VALUES (@email, @role, @company, @salary, @type, @location)`

	args := pgx.NamedArgs{
		"email":    form.Email,
		"role":     form.Role,
		"company":  form.Company,
		"salary":   form.Salary,
		"type":     form.Type,
		"location": form.Location,
	}

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to create vacancy")
		return err
	}
	return nil
}

func (r *PostgresRepository) GetVacancies(ctx context.Context, limit int, offset int) ([]Vacancy, error) {
	query := `SELECT * FROM vacancies ORDER BY created_at DESC LIMIT @limit OFFSET @offset`
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}
	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to get vacancies")
		return nil, err
	}
	defer rows.Close()

	vacancies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Vacancy])
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to collect rows")
		return nil, err
	}
	return vacancies, nil
}
