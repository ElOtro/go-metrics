package postgres

type pgstorage struct {
}

func New() *pgstorage {
	return &pgstorage{}
}

func (m pgstorage) Get() string {
	return ""
}

func (m pgstorage) Set(t, n, v string) error {
	return nil
}
