package com.spaceuptech.dbevents.database

import akka.actor.typed.{ActorRef, Behavior}
import akka.actor.typed.scaladsl.Behaviors
import com.spaceuptech.dbevents.spacecloud.{DatabaseConfig, EventsSink}

object Database {
  def createActor(projectId: String, dbType: String, actor: ActorRef[EventsSink.Command]): Behavior[Command] = {
    dbType match {
      case "postgres" | "mysql" | "sqlserver" => Behaviors.withTimers[Command](timers => Behaviors.setup[Command](context => new Debezium(context, timers, projectId, actor)))
      case "mongo" => Behaviors.withTimers[Command](timers => Behaviors.setup[Command](context => new Mongo(context, timers, projectId, actor)))
      case _ => throw new Exception(s"Invalid db type ($dbType) provided")
    }
  }

  sealed trait Command
  case class ChangeRecord(payload: ChangeRecordPayload, project: String, dbAlias: String, dbType: String) extends Command
  case class CheckEngineStatus() extends Command
  case class UpdateEngineConfig(config: DatabaseConfig) extends Command
  case class ProcessEngineConfig(conn: String, config: DatabaseConfig) extends Command
  case class Stop() extends Command
}

