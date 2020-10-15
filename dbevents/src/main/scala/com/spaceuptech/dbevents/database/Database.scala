package com.spaceuptech.dbevents.database

import akka.actor.typed.Behavior
import akka.actor.typed.scaladsl.Behaviors
import com.spaceuptech.dbevents.spacecloud.DatabaseConfig

object Database {
  def createActor(projectId: String, config: DatabaseConfig): Behavior[Command] = {
    config.`type` match {
      case "postgres" | "mysql" | "sqlserver" => Behaviors.withTimers[Command](timers => Behaviors.setup[Command](context => new Debezium(context, timers, projectId, config)))
      case _ => throw new Exception(s"Invalid db type (${config.`type`}) provided")
    }
  }

  sealed trait Command
  case class CheckEngineStatus() extends Command
  case class Stop() extends Command
}

