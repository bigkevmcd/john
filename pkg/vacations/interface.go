package vacations

type VacationService interface {
	ActiveVacation(string) (*Vacation, error)
}

type SMTPClient interface {
}
