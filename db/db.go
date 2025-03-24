package db

import (
	gorestdb "github.com/cuityhj/gorest/db"
	gorestresource "github.com/cuityhj/gorest/resource"
)

const (
	ConnStrTemplate = "user=%s password=%s host=%s port=%d database=%s sslmode=disable"
)

func NewDBConn(driverName gorestdb.DriverName, resources []gorestresource.Resource, connStr string, dropSchemaList ...string) (gorestdb.ResourceStore, error) {
	meta, err := gorestdb.NewResourceMeta(resources)
	if err != nil {
		return nil, err
	}

	if driverName == "" {
		driverName = gorestdb.DriverNamePostgresql
	}

	return gorestdb.NewRStore(driverName, connStr, meta, dropSchemaList...)
}
