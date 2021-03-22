package main

import (
	"fmt"
)

func parseConnectionString(dbType, conn string) error {
	switch dbType {
	case "mysql":
		return parseMySQLConn(conn)
	case "postgres":
		return parsePostgresConnString(conn)
	case "sqlserver":
		return parseSQLSeverConnString(conn)
	default:
		return fmt.Errorf("invalid dbtype (%s) provided", dbType)
	}
}
