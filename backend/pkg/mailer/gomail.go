package mailer

type Mailer interface {
	Send(to []string, subject string, body string, isHTML bool) error
}
