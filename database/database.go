package database

import (
	"github.com/gobuffalo/pop"
	"fmt"
)

type Database struct
{
	connection *pop.Connection
	healthCheck func (connection *pop.Connection) error
}


func (d *Database) Connection() (*pop.Connection){
	return d.connection
}

func (d *Database) Open() error {
	if err := d.connection.Open(); err != nil {
		return err
	}
	return d.HealthCheck()
}

func (d *Database) Close() error {
	return d.connection.Close()
}

func (d *Database) HealthCheck() error {
	return d.healthCheck(d.connection)
}

func New(cd *pop.ConnectionDetails, healthCheck func (connection *pop.Connection)error) (*Database, error) {
	connection, err := pop.NewConnection(cd)
	if err != nil {
		return nil, fmt.Errorf("failed to create the database connection %s", err)
	}
	return &Database{connection: connection, healthCheck: healthCheck}, nil
}






