package main

type infra struct {
	db     DBAccess
	logger Logger
}

func NewInfraHandler(db DBAccess, logger Logger) *infra {
	i := &infra{db: db, logger: logger}
	return i
}
func (i *infra) TryMarkMessageAsProcessed(messageId string) (bool, error) {
	alreadyProcessed, err := i.db.TryMarkMessageAsProcessed(messageId)
	return alreadyProcessed, err
}
