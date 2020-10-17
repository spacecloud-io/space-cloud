package com.spaceuptech.dbevents.database

import akka.actor.typed.{ActorSystem, Behavior}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, TimerScheduler}
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

  // Make the mongo db driver
  val connectionString = new MongoClientURI(connString)
  val mongoClient = new MongoClient(connectionString)
  override def onMessage(msg: Database.Command): Behavior[Database.Command] = {
    msg match {
      case CheckEngineStatus() =>
        this
    }
  }
}
