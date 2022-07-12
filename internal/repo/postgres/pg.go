package postgres

type pgstorage struct {
}

func New() *pgstorage {
	return &pgstorage{}
}

func (m pgstorage) Get(t, n string) (string, error) {
	return "", nil
}

func (m pgstorage) GetAll() string {
	return ""
}

func (m pgstorage) Set(t, n, v string) error {
	return nil
}
