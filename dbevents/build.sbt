name := "db-events-soruce"

version := "0.1.0"

scalaVersion := "2.13.1"

lazy val akkaVersion = "2.6.9"

libraryDependencies ++= Seq(
  "com.typesafe.akka" %% "akka-actor-typed" % akkaVersion,
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

// Dependencies for parsing JSON
libraryDependencies += "org.json4s" %% "json4s-jackson" % "3.7.0-M6"
