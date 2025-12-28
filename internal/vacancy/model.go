package vacancy

import "time"

type Vacancy struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Role      string    `db:"role"`
	Company   string    `db:"company"`
	Salary    string    `db:"salary"`
	Type      string    `db:"type"`
	Location  string    `db:"location"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
