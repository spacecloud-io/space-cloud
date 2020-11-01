package com.spaceuptech.dbevents.database

import java.util.concurrent.Executors

import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import akka.actor.typed.{ActorRef, ActorSystem, Behavior, PostStop, Signal}
import com.spaceuptech.dbevents.pubsub.RabbitMQ
import com.spaceuptech.dbevents.spacecloud.DatabaseConfig

import scala.concurrent.ExecutionContextExecutor
import scala.concurrent.duration.DurationInt
import scala.util.{Failure, Success}

class Mongo(context: ActorContext[Database.Command], timers: TimerScheduler[Database.Command], projectId: String, broker: ActorRef[RabbitMQ.Command]) extends AbstractBehavior[Database.Command](context) {

  import Database._

  // Lets get the connection string first
  implicit val system: ActorSystem[Nothing] = context.system
  implicit val executionContext: ExecutionContextExecutor = system.executionContext
  private val executor = Executors.newSingleThreadExecutor

  // The status variables
  private var name: String = ""
  private var connString: String = ""
  private var status: Option[MongoStatus] = None
  private var dbConfig: DatabaseConfig = _

  // Start task for status check
  timers.startTimerAtFixedRate(CheckEngineStatus(), 30.second)

  override def onMessage(msg: Database.Command): Behavior[Database.Command] = {
    msg match {
      case CheckEngineStatus() =>
        status match {
          case Some(value) =>
            if (value.future.isCancelled || value.future.isDone) {
              // Just making sure its closed first
              value.future.cancel(true)
              value.store.stop()

              println(s"Mongo watcher $name is closed. Restarting...")
              status = Some(Utils.startMongoWatcher(projectId, dbConfig.dbAlias, connString, dbConfig.name, executor, context.self))
            }

          // No need to do anything if status isn't defined
          case None =>
        }
        this

      case record: ChangeRecord =>
        broker ! RabbitMQ.EmitEvent(record)
        this

      case UpdateEngineConfig(config) =>
        // Start mongo watcher
        getConnString(projectId, config.conn) onComplete {
          case Success(conn) =>
            // Simply return if there are no changes to the connection string
            if (conn == connString) return this

            // Store the connection string for future reference
            this.connString = conn

            // Kill the previous debezium engine
            stopOperations()

            // Store the name and db config object for later use
            name = s"$projectId:${config.dbAlias}"
            dbConfig = config

            // Start the engine
            println(s"Staring debezium engine $name")

            status = Some(Utils.startMongoWatcher(projectId, config.dbAlias, connString, config.name, executor, context.self))

          case Failure(ex) =>
            println(s"Unable to get connection string for mongo engine ($name) - ${ex.getMessage}")

        }
        this

      case Stop() => Behaviors.stopped

    }
  }

  private def stopOperations(): Unit = {
    // Shut down the debezium engine
    status match {
      case Some(value) =>
        if (!value.future.isCancelled && !value.future.isDone) {
          println(s"Closing mongo watcher - $name")
          value.future.cancel(true)
          value.store.stop()
          println(s"Closed mongo watcher - $name")
        }

      // No need to do anything if status isn't defined
      case None =>
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      // Shutdown the timer
      timers.cancelAll()

      // Stop engine operations
      stopOperations()

      // Shut down the main executor
      executor.shutdownNow()

      this
  }
}
