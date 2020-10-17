package com.spaceuptech.dbevents.database

import java.util.Properties
import java.util.concurrent.{Callable, ExecutorService}

import akka.actor.typed.ActorRef
import com.mongodb.client.model.changestream.ChangeStreamDocument
import com.mongodb.{Block, MongoClient, MongoClientURI}
import com.spaceuptech.dbevents.database.Database.ChangeRecord
import com.spaceuptech.dbevents.{DatabaseSource, Global}
import io.debezium.engine.format.Json
import io.debezium.engine.DebeziumEngine
import org.bson.Document
import org.json4s._
import org.json4s.jackson.JsonMethods._

object Utils {
  def startMongoWatcher(projectId: String, conn: String, dbName: String, executorService: ExecutorService, actor: ActorRef[Database.Command]): java.util.concurrent.Future[_] = {
    val connectionString = new MongoClientURI("mongodb://localhost:27017")
    val mongoClient = new MongoClient(connectionString)
    val db = mongoClient.getDatabase(dbName)

    val printBlock = new Block[ChangeStreamDocument[Document]]() {
      override def apply(t: ChangeStreamDocument[Document]): Unit = {

      }
    }

    executorService.submit(new Callable[Unit] {
      override def call(): Unit = {
        val itr = db.watch().iterator()

      }
    })
  }

  def startDebeziumEngine(source: DatabaseSource, executorService: ExecutorService, actor: ActorRef[Database.Command]): DebeziumStatus = {
    // Create the engine configuration object
    val props = source.dbType match {
      case "mysql" => prepareMySQLConfig(source)
      case "postgres" => preparePostgresConfig(source)
      case "sqlserver" => prepareSQLServerConfig(source)
      case _ => throw new IllegalArgumentException
    }

    // Create a new debezium engine

    val engine = DebeziumEngine.create(classOf[Json]).using(props).notifying(record => {
      implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
      // Extract the change feed value
      val jsonString = record.value()

      // Marshal the string only if the json string is not null
      if (jsonString != null) {
        // Parse the json value and forward it to our actor
        val payload = parse(jsonString).extract[ChangeRecordPayload]
        actor ! ChangeRecord(payload, source.project, source.dbAlias, source.dbType)
      }

    }).build()

    // Run the engine asynchronously
    val future = executorService.submit(engine)
    DebeziumStatus("", future, engine)
  }

  def prepareMySQLConfig(source: DatabaseSource): Properties = {
    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", generateConnectorName(source))
    props.setProperty("connector.class", "io.debezium.connector.mysql.MySqlConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
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
    props.setProperty("database.history", getDatabaseHistoryStorageClass)
    props.setProperty("database.history.file.filename", "./dbhistory.dat")

    props
  }

  def prepareSQLServerConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.postgresql.PostgresConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
    props.setProperty("offset.storage.file.filename", "./offsets.dat")
    props.setProperty("offset.flush.interval.ms", "60000")
    props.setProperty("converter.schemas.enable", "false")
    /* begin connector properties */
    props.setProperty("database.hostname", source.config.getOrElse("host", "localhost"))
    props.setProperty("database.port", source.config.getOrElse("port", "1433"))
    props.setProperty("database.user", source.config.getOrElse("user", "sa"))
    props.setProperty("database.password", source.config.getOrElse("password", "mypassword"))
    props.setProperty("database.dbname", source.config.getOrElse("db", "test"))
    props.setProperty("database.server.name", s"${generateConnectorName(source)}_connector")

    props
  }



  def preparePostgresConfig(source: DatabaseSource): Properties = {
    val name = generateConnectorName(source)

    val props = io.debezium.config.Configuration.empty().asProperties()
    props.setProperty("name", name)
    props.setProperty("connector.class", "io.debezium.connector.postgresql.PostgresConnector")
    props.setProperty("offset.storage", getOffsetStorageClass)
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
    props.setProperty("database.sslmode", source.config.getOrElse("sslMode", "disable"))

    props
  }

  def getOffsetStorageClass: String = {
    Global.storageType match {
      case "k8s" => "com.spaceuptech.dbevents.database.KubeOffsetBackingStore"
      case _ => "org.apache.kafka.connect.storage.FileOffsetBackingStore"
    }
  }

  def getDatabaseHistoryStorageClass: String = {
    Global.storageType match {
      case "k8s" => "com.spaceuptech.dbevents.database.KubeDatabaseHistory"
      case _ => "io.debezium.relational.history.FileDatabaseHistory"
    }
  }

  def generateConnectorName(source: DatabaseSource): String = {
    s"${source.project}_${source.dbAlias}"
  }
}
