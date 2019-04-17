package config

/**
 * @author ollyel, ollykel416@gmail.com
 * @date Apr 22, 2019
 * Parses connection string for databases from ConnConfig struct
 * Uses github.com/go-sql-driver/mysql for its Config struct, DSN parsing
 */

import (
	"github.com/go-sql-driver/mysql"
)

func ParseConnMySQL (cfg *ConnConfig) (_ string, err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.User
	mysqlConfig.Net = cfg.Protocol
	mysqlConfig.Addr = cfg.Host + ":" + cfg.Port
	mysqlConfig.DBName = cfg.DBName
	mysqlConfig.Params = cfg.Params
	mysqlConfig.ServerPubKey = cfg.SSLKey
	mysqlConfig.Passwd, err = GetPassword(cfg.Auth)
	if err != nil { return "", err }
	return mysqlConfig.FormatDSN(), nil
}//-- end func ParseConnMySQL

func ParseConnPostgres (cfg *ConnConfig) (_ string, err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.User
	mysqlConfig.Net = cfg.Protocol
	mysqlConfig.Addr = cfg.Host + ":" + cfg.Port
	mysqlConfig.DBName = cfg.DBName
	mysqlConfig.Params = cfg.Params
	// ssl mode
	if mysqlConfig.Params == nil {
		mysqlConfig.Params = make(map[string]string)
	}
	mysqlConfig.Params["sslmode"] = cfg.SSLMode.String()
	mysqlConfig.ServerPubKey = cfg.SSLKey
	mysqlConfig.Passwd, err = GetPassword(cfg.Auth)
	if err != nil { return "", err }
	return "postgres://" + mysqlConfig.FormatDSN(), nil
}//-- end func ParseConnMySQL

func ParseConnMongo (cfg *ConnConfig) (_ string, err error) {
	mysqlConfig := mysql.NewConfig()
	mysqlConfig.User = cfg.User
	// Protocol (i.e. tcp) not used for mongo
	mysqlConfig.Net = cfg.Host + ":" + cfg.Port
	mysqlConfig.DBName = cfg.DBName
	mysqlConfig.Params = cfg.Params
	if mysqlConfig.Params == nil {
		mysqlConfig.Params = make(map[string]string)
	}
	// ssl mode
	if cfg.SSLMode == Disable {
		mysqlConfig.Params["ssl"] = "false"
	} else { mysqlConfig.Params["ssl"] = "true" }
	mysqlConfig.ServerPubKey = cfg.SSLKey
	mysqlConfig.Passwd, err = GetPassword(cfg.Auth)
	if err != nil { return "", err }
	return "mongodb://" + mysqlConfig.FormatDSN(), nil
}//-- end func ParseConnMongo

