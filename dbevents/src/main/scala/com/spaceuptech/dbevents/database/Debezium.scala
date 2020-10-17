package com.spaceuptech.dbevents.database

import java.util.concurrent.Executors

import akka.actor.typed.{ActorSystem, Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import com.spaceuptech.dbevents.{DatabaseSource, Global}
import com.spaceuptech.dbevents.spacecloud._
import org.json4s.DefaultFormats
import org.json4s.jackson.JsonMethods.parse

import scala.concurrent.{Await, ExecutionContextExecutor, Future}
import scala.concurrent.duration.DurationInt
import sys.process._

class Debezium(context: ActorContext[Database.Command], timers: TimerScheduler[Database.Command], projectId: String, config: DatabaseConfig) extends AbstractBehavior[Database.Command](context) {
  import Database._

  // Lets get the connection string first
  implicit val system: ActorSystem[Nothing] = context.system
  implicit val executionContext: ExecutionContextExecutor = system.executionContext
  private val connString = Await.result(getConnString(projectId, config.conn), 10.seconds)
  private val source = generateDatabaseSource(projectId, connString, config)

  // Extract name of actor
  private val name = Utils.generateConnectorName(source)
  context.log.info(s"Staring debezium engine $name")

  // Start the debezium engine
  private val executor = Executors.newSingleThreadExecutor
  private var status = Utils.startDebeziumEngine(source, executor, context.self)

  // Start task for status check
  timers.startTimerAtFixedRate(name, CheckEngineStatus(), 30.second)

  override def onMessage(msg: Command): Behavior[Command] = {
    // No need to handle any messages
    msg match {
      case Database.CheckEngineStatus() =>
        // Try starting the debezium engine again only if it wasn't running already
        if (status.future.isDone || status.future.isCancelled) {
          // Just making sure its closed first
          status.engine.close()
          status.future.cancel(true)

          context.log.info(s"Debezium engine $name is closed. Restarting...")
          status = Utils.startDebeziumEngine(source, executor, context.self)
        }
        this

      case ChangeRecord(payload, project, dbAlias, dbType) =>
        this

      case Stop() => Behaviors.stopped
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      // Shutdown the timer
      timers.cancelAll()

      // Shut down the debezium engine
      if (!status.future.isCancelled && !status.future.isDone) {
        context.log.info(s"Closing debezium engine - $name")
        status.engine.close()
        status.future.cancel(true)
        context.log.info(s"Closed debezium engine - $name")
      }

      // Shut down the main executor
      executor.shutdownNow()

      this
  }

  private def generateDatabaseSource(projectId: String, conn: String, db: DatabaseConfig): DatabaseSource = {
    // Generate ast from the database conn string
    val jsonString = s"conn-string-parser parse --db-type ${db.`type`} '$conn'".!!

    // Parse the ast into config object
    implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
    var config = parse(jsonString).extract[Map[String, String]]

    // Make necessary adjustments
    db.`type` match {
      case "mysql" => config += "db" -> db.name
      case "postgres" => config += "schema" -> db.name
    }

    DatabaseSource(projectId, db.dbAlias, db.`type`, config)
  }

}
