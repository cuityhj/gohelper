package postgresql

import (
	gorestdb "github.com/zdnscloud/gorest/db"
	gorestresource "github.com/zdnscloud/gorest/resource"
)

const ConnStrTemplate = "user=%s password=%s host=%s port=%d database=%s sslmode=disable pool_max_conns=10"

func NewDBConn(resources []gorestresource.Resource, connStr string) (gorestdb.ResourceStore, error) {
	meta, err := gorestdb.NewResourceMeta(resources)
	if err != nil {
		return nil, err
	}

	return gorestdb.NewRStore(connStr, meta)
}
