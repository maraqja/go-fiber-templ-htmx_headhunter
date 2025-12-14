package vacancy

type VacancyCreateForm struct {
	Role     string `form:"role"`
	Company  string `form:"company"`
	Salary   string `form:"salary"`
	Type     string `form:"type"`
	Location string `form:"location"`
	Email    string `form:"email"`
}
