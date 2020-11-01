#!/bin/sh
# NOTE: Run this script from the gateway directory
# 1) ensure port 4122,3306,5432,1433,27017 are not busy
# 2) ensure docker & golang is installed
set -e

echo "building and running gateway binary"
go build .

# run the gateway in background
sudo kill -9 `sudo lsof -t -i:4122` &
rm config.yaml &
./gateway run --dev --log-format text &

sleep 10
#sc_process_num=$!

# change directory gateway sql package
cd modules/crud/sql

## mysql test
echo "starting mysql container, it will take 30 seconds"
docker run --name integration-mysql -p 3306:3306  -e MYSQL_ROOT_PASSWORD=my-secret-pw -d mysql:8.0.21 &> /dev/null
sleep 30
echo "running integration tests for mysql"
go test -tags integration -db_type mysql -conn "root:my-secret-pw@tcp(localhost:3306)/"
echo "removing mysql container"
docker rm -f integration-mysql
echo "\n"

# postgres test
echo "starting postgres container, it will take 30 seconds"
docker run --name integration-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:13
sleep 30
echo "running integration tests for postgres"
go test -tags integration -db_type postgres -conn "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
echo "removing postgres container"
docker rm -f integration-postgres
echo "\n"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

echo "starting postgres container, it will take 30 seconds"
docker run --name integration-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:12.4
sleep 30
echo "running integration tests for postgres"
go test -tags integration -db_type postgres -conn "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
echo "removing postgres container"
docker rm -f integration-postgres
echo "\n"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

echo "starting postgres container, it will take 30 seconds"
docker run --name integration-postgres -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres:11.9
sleep 30
echo "running integration tests for postgres"
go test -tags integration -db_type postgres -conn "postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable"
echo "removing postgres container"
docker rm -f integration-postgres
echo "\n"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

## sqlserver test
echo "starting sqlserver container,it will take 30 seconds"
docker run --name integration-sqlserver -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=yourPassword@#12345' -p 1433:1433 -d mcr.microsoft.com/mssql/server:2019-latest
sleep 30
echo "running integration tests for sqlserver"
go test -tags integration -db_type sqlserver -conn "Data Source=localhost,1433;Initial Catalog=master;User ID=sa;Password=yourPassword@#12345;"
docker rm -f integration-sqlserver
echo "removing sqlserver container"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

cd ../mgo

# mongo test
echo "starting mongo container,it will take 30 seconds"
docker run --name integration-mongo -p 27017:27017 -d mongo:4.4
sleep 30
echo "running integration tests for mongo"
go test -tags integration -db_type mongo -conn "mongodb://localhost:27017"
echo "removing mongo container"
docker rm -f integration-mongo
echo "\n"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

echo "starting mongo container,it will take 30 seconds"
docker run --name integration-mongo -p 27017:27017 -d mongo:4.2
sleep 30
echo "running integration tests for mongo"
go test -tags integration -db_type mongo -conn "mongodb://localhost:27017"
echo "removing mongo container"
docker rm -f integration-mongo
echo "\n"

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
../../../gateway run --dev --log-format text &

echo "starting mongo container,it will take 30 seconds"
docker run --name integration-mongo -p 27017:27017 -d mongo:4.0
sleep 30
echo "running integration tests for mongo"
go test -tags integration -db_type mongo -conn "mongodb://localhost:27017"
echo "removing mongo container"
docker rm -f integration-mongo

sudo kill -9 `sudo lsof -t -i:4122`
rm ../../../config.yaml &
