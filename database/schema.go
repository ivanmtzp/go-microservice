package database

import (
	"github.com/gobuffalo/pop"
)


type Schema struct {
	properties *Properties
	connection *pop.Connection
}

func NewSchema(p *Properties) (*Schema, error) {
	connection, err := createConnection(p)
	if err != nil {
		return nil, err
	}
	schema := &Schema{
		properties: p,
		connection: connection,
	}

	return schema, nil
}

func (s *Schema) CreateDatabase() error {
	return pop.CreateDB(s.connection)
}

func (s *Schema) DropDatabase() error {
	return pop.DropDB(s.connection)
}

func (s* Schema) MigrateDatabase(migrationsPath string) error {
	mig, err := pop.NewFileMigrator(migrationsPath, s.connection)
	if err != nil {
		return err
	}
	return mig.Up()
}

func (s *Schema) RecreateDatabase(migrationsPath string) error {
	s.DropDatabase()
	if err := s.CreateDatabase(); err != nil {
		return err
	}
	return s.MigrateDatabase(migrationsPath)
}

