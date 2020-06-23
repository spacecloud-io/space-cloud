# NOTE: Run this script from the gateway directory
/bin/sh

echo "building and running gateway binary"
go build .
./gateway run --dev

# change directory gateway sql package
cd modules/crud/sql

# mysql test
echo "starting mysql container"
docker run --name integration-mysql -p 3306:3306  -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:latest
echo "running integration tests for mysql"
go test -tags integration -db_type mysql -conn "root:my-secret-pw@tcp(localhost:3306)/myproject"
echo "removing mysql container"
docker rm -f integration-mysql

# postgres test
echo "starting postgres container"
docker run --name integration-postgres -p 5432:5432 -e POSTGRES_PASSWORD=my-secret-pw -d postgres
echo "running integration tests for postgres"
go test -tags integration -db_type postgres -conn "postgres://postgres:my-secret-pw@localhost:5432/postgres?sslmode=disable"
echo "removing postgres container"
docker rm -f integration-postgres

# sqlserver test
echo "starting sqlserver container"
docker run --name integration-sqlserver -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=my-secret-pw' -p 1433:1433 -d mcr.microsoft.com/mssql/server:2017-CU8-ubuntu
echo "running integration tests for sqlserver"
go test -tags integration -db_type sqlserver -conn "Data Source=localhost,1433;Initial Catalog=master;User ID=sa;Password=my-secret-pw;"
docker rm -f integration-sqlserver
echo "removing sqlserver container"

# mongo test
# switch to mongo directory
cd ../mgo
echo "starting mongo container"
docker run --name integration-mongo -p 27017:27017 -d mongo:latest
echo "running integration tests for mongo"
go test -tags integration -db_type mongo -conn "mongodb://localhost:27017"
docker rm -f integration-mongo
echo "removing mongo container"