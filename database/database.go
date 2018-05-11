package database

import (
	"github.com/gobuffalo/pop"
	"fmt"
	"strconv"
)

type Database struct
{
	properties *Properties
	connection *pop.Connection
	healthCheckQuery string
}

type Properties struct {
	Dialect string
	Database string
	Host string
	Port int
	User string
	Password string
	Pool int
}

func (d *Database) Address() string {
	return fmt.Sprintf("%s:%d", d.properties.Host, d.properties.Port)
}

func (d *Database) Connection() (*pop.Connection){
	return d.connection
}

func (d *Database) Close() error {
	if d.connection != nil {
		return d.connection.Close()
	}
	return nil
}

func createConnection(p *Properties) (*pop.Connection, error) {
	cd := &pop.ConnectionDetails{
		Dialect: p.Dialect,
		Database: p.Database,
		Host: p.Host,
		Port: strconv.Itoa(p.Port),
		User: p.User,
		Password: p.Password,
		Pool: p.Pool,
		IdlePool: 0,
	}
	connection, err := pop.NewConnection(cd)
	if err != nil {
		return nil, fmt.Errorf("failed to create the database connection %s", err)
	}
	if err := connection.Open(); err != nil {
		return nil, fmt.Errorf("failed to open database connection %s", err)
	}
	return connection, nil
}

func (d *Database) HealthCheck() error {
	var response []interface{}
	return d.connection.RawQuery(d.healthCheckQuery).All(&response)
}

func NewDatabase(p *Properties, healthCheckQuery string) (*Database, error) {
	connection, err := createConnection(p)
	if err != nil {
		return nil, err
	}
	db := &Database{
		properties: p,
		connection: connection,
		healthCheckQuery: healthCheckQuery}
	if err := db.HealthCheck(); err != nil {
		return nil, fmt.Errorf("failed database healthcheck %s", err)
	}
	return db, nil
}

