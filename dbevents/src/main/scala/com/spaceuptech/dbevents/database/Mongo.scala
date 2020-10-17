package com.spaceuptech.dbevents.database

import java.util.concurrent.Executors

import akka.actor.typed.{ActorSystem, Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import com.mongodb.{MongoClient, MongoClientURI}
import com.spaceuptech.dbevents.spacecloud.DatabaseConfig

import scala.concurrent.duration.DurationInt
import scala.concurrent.{Await, ExecutionContextExecutor}

class Mongo(context: ActorContext[Database.Command], timers: TimerScheduler[Database.Command], projectId: String, config: DatabaseConfig) extends AbstractBehavior[Database.Command](context) {
  import Database._

  // Lets get the connection string first
  implicit val system: ActorSystem[Nothing] = context.system
  implicit val executionContext: ExecutionContextExecutor = system.executionContext
  private val connString = Await.result(getConnString(projectId, config.conn), 10.seconds)

  // Start mongo watcher
  private val executor = Executors.newSingleThreadExecutor
  private var status = Utils.startMongoWatcher(projectId, config.dbAlias, connString, config.name, executor, context.self)

  // Start task for status check
  timers.startTimerAtFixedRate(s"$projectId:${config.dbAlias}", CheckEngineStatus(), 30.second)

  override def onMessage(msg: Database.Command): Behavior[Database.Command] = {
    msg match {
      case CheckEngineStatus() =>
        if (status.isCancelled || status.isDone) {
          // Just making sure its closed first
          status.cancel(true)

          context.log.info(s"Mongo watcher $projectId:${config.dbAlias} is closed. Restarting...")
          status = Utils.startMongoWatcher(projectId, config.dbAlias, connString, config.name, executor, context.self)
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
      if (!status.isCancelled && !status.isDone) {
        context.log.info(s"Closing mongo watcher - $projectId:${config.dbAlias}")
        status.cancel(true)
        context.log.info(s"Closed mongo watcher - $projectId:${config.dbAlias}")
      }

      // Shut down the main executor
      executor.shutdownNow()

      this
  }
}
