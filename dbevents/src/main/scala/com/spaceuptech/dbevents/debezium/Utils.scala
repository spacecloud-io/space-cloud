package com.spaceuptech.dbevents.debezium

import java.util.Properties
import java.util.concurrent.{CompletableFuture, ExecutorService, Executors, Future, LinkedBlockingQueue, ThreadPoolExecutor, TimeUnit}

import com.spaceuptech.dbevents.DatabaseSource
import io.debezium.engine.format.{ChangeEventFormat, Json}
import io.debezium.engine.{ChangeEvent, DebeziumEngine}
import org.json4s._
import org.json4s.jackson.JsonMethods._

object Utils {

  case class ChangeRecord(payload: ChangeRecordPayload, project: String, dbAlias: String, dbType: String)

  case class ChangeRecordPayload(op: Option[String], before: Option[Map[String, Any]], after: Option[Map[String, Any]], source: ChangeRecordPayloadSource)

  case class ChangeRecordPayloadSource(name: String, ts_ms: Long, table: String)

  case class DebeziumStatus(future: Future[_], engine: DebeziumEngine[ChangeEvent[String, String]])

  def startDebeziumEngine(source: DatabaseSource, executorService: ExecutorService): DebeziumStatus = {
    // Create the engine configuration object
    val props = source.dbType match {
      case "mysql" => prepareMySQLConfig(source)
      case "postgres" => preparePostgresConfig(source)
      case _ => throw new IllegalArgumentException
    }

    // Create a new debezium engine

    val engine = DebeziumEngine.create(classOf[Json]).using(props).notifying((record) => {
      implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
      // Extract the change feed value
      val jsonString = record.value()

      // Marshal the string only if the json string is not null
      if (jsonString != null) {
        // Parse the json value

        System.out.println()
        System.out.println("DB Record Raw:", jsonString)
        val payload = parse(jsonString).extract[ChangeRecordPayload]
        System.out.println("DB Record:", ChangeRecord(payload, source.project, source.dbAlias, source.dbType))
        System.out.println()
      }

    }).build()

    // Run the engine asynchronously

    val future = executorService.submit(engine)

    DebeziumStatus(future, engine)
  }

  def prepareMySQLConfig(source: DatabaseSource): Properties = {
    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", generateConnectorName(source))
    props.setProperty("connector.class", "io.debezium.connector.mysql.MySqlConnector")
    props.setProperty("offset.storage", "org.apache.kafka.connect.storage.FileOffsetBackingStore")
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "3306"))
    props.setProperty("database.user", source.config.getOrElse("user", "root"))
    props.setProperty("database.password", source.config.getOrElse("password", "my-secret-pw"))
    props.setProperty("database.include.list", source.config.getOrElse("db", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")
    props.setProperty("database.ssl.mode", source.config.getOrElse("sslMode", "disabled"))
    props.setProperty("database.history", "io.debezium.relational.history.FileDatabaseHistory")
    props.setProperty("database.history.file.filename", "./dbhistory.dat")

    props
  }

  def preparePostgresConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.postgresql.PostgresConnector")
    props.setProperty("offset.storage", "org.apache.kafka.connect.storage.FileOffsetBackingStore")
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("plugin.name", "pgoutput")
    props.setProperty("slot.name", name)
    props.setProperty("publication.name", name)
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "5432"))
    props.setProperty("database.user", source.config.getOrElse("user", "postgres"))
    props.setProperty("database.password", source.config.getOrElse("password", "mysecretpassword"))
    props.setProperty("database.dbname", source.config.getOrElse("db", "postgres"))
    props.setProperty("schema.include.list", source.config.getOrElse("schema", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")
    props.setProperty("database.ssl.mode", source.config.getOrElse("sslMode", "disabled"))

    props
  }

  def prepareMongoConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.mongodb.MongoDbConnector")
    props.setProperty("offset.storage", "org.apache.kafka.connect.storage.FileOffsetBackingStore")
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("mongodb.hosts", source.config.getOrElse("hosts", "localhost:27017"))
    props.setProperty("mongodb.name", source.config.getOrElse("name", name))
    props.setProperty("mongodb.user", source.config.getOrElse("user", "user"))
    props.setProperty("mongodb.password", source.config.getOrElse("password", "pass"))
    props.setProperty("mongodb.authsource", source.config.getOrElse("authSource", "admin"))
    props.setProperty("mongodb.ssl.enabled",source.config.getOrElse("sslEnabled", "false"))
    props.setProperty("database.include.list", source.config.getOrElse("db", "test"))

    props
  }

  def generateConnectorName(source: DatabaseSource): String = {
    s"${source.project}_${source.dbAlias}"
  }
}
