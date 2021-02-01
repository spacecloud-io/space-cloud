name := "db-events-soruce"

version := "0.2.0"

scalaVersion := "2.13.1"

lazy val akkaVersion = "2.6.9"
lazy val akkaHttpVersion = "10.2.1"

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-actor-typed" % akkaVersion,
  "com.typesafe.akka" %% "akka-stream" % akkaVersion,
  "com.typesafe.akka" %% "akka-http" % akkaHttpVersion,
  "ch.qos.logback" % "logback-classic" % "1.2.3",
  "com.typesafe.akka" %% "akka-actor-testkit-typed" % akkaVersion % Test,
  "org.scalatest" %% "scalatest" % "3.1.0" % Test
)

// Dependencies for Debezium
libraryDependencies ++= Seq(
  "io.debezium" % "debezium-embedded" % "1.3.0.Final",
  "io.debezium" % "debezium-api" % "1.3.0.Final",

  // Database connectors
  "io.debezium" % "debezium-connector-mongodb" % "1.3.0.Final",
  "io.debezium" % "debezium-connector-mysql" % "1.3.0.Final",
  "io.debezium" % "debezium-connector-sqlserver" % "1.3.0.Final",
  "io.debezium" % "debezium-connector-postgres" % "1.3.0.Final"
)

// Dependencies for the mongo db driver
libraryDependencies += "org.mongodb" % "mongodb-driver" % "3.12.7"


// Dependencies for parsing JSON
libraryDependencies += "org.json4s" %% "json4s-jackson" % "3.7.0-M6"

// Dependencies to play around with JWT
libraryDependencies += "com.auth0" % "java-jwt" % "3.11.0"

// Dependencies for the kubernetes client
libraryDependencies += "io.kubernetes" % "client-java" % "10.0.0"

// https://mvnrepository.com/artifact/commons-io/commons-io
libraryDependencies += "commons-io" % "commons-io" % "2.8.0"

libraryDependencies += "com.typesafe" % "config" % "1.4.1"

enablePlugins(JavaAppPackaging)
