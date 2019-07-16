export const dbTypes = {
  MONGO: "mongo",
  POSTGRESQL: "sql-postgres",
  MYSQL: "sql-mysql"
}

export const defaultDbConnectionStrings = {
  [dbTypes.MONGO]: "mongodb://localhost:27017",
  [dbTypes.POSTGRESQL]: "postgres://postgres:mysecretpassword@localhost/postgres?sslmode=disable",
  [dbTypes.MYSQL]: "user:my-secret-pwd@/project"
}

export const PAYU_MERCHANT_KEY = "FkS68IbH"

export const SPACE_API_PROJECT = "space-cloud"
export const SPACE_API_URL = "http://localhost:8080"

