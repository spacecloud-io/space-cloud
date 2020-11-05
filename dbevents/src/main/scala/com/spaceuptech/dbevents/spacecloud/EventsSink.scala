package com.spaceuptech.dbevents.spacecloud

import java.time.{Instant, OffsetDateTime, ZoneId}
import java.util.Calendar

import akka.actor.typed.{ActorSystem, Behavior}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.spaceuptech.dbevents.database.Database.{ChangeRecord, Command}

import scala.concurrent.ExecutionContextExecutor
import scala.util.{Failure, Success}

object EventsSink {

  def apply(projectId: String): Behavior[Command] = Behaviors.setup(context => new EventsSink(context, projectId))

  sealed trait Command
  case class EmitEvent(record: ChangeRecord) extends Command
}

class EventsSink(context: ActorContext[EventsSink.Command], projectId: String) extends AbstractBehavior[EventsSink.Command](context) {
  implicit val system: ActorSystem[Nothing] = context.system
  implicit val executionContext: ExecutionContextExecutor = system.executionContext

  override def onMessage(msg: EventsSink.Command): Behavior[EventsSink.Command] = {
    msg match {
      case EventsSink.EmitEvent(record) =>
        // Queue event in gateway
        queueEvent(projectId, prepareQueueRequest(record)).onComplete {
          case Success(_) => println("Event logged successfully")
          case Failure(exception) => println(s"Unable to log event - ${exception.getMessage}")
        }
        this
    }
  }

  private def prepareQueueRequest(record: ChangeRecord): QueueEvent = {
    QueueEvent(
      `type` = record.payload.op match {
        case "c" => "DB_INSERT"
        case "u" => "DB_UPDATE"
        case "d" => "DB_DELETE"
        case _ => "UNKNOWN_OP"
      },
      payload = DatabaseEvent(
        db = record.dbAlias,
        col = record.payload.source.table,
        doc = record.payload.after,
        find = record.payload.before,
        before = record.dbType match {
          case "mongo" => None
          case _ => record.payload.before
        }
      ),
      options = DatabaseEventOptions(
        db = record.dbAlias,
        col = record.payload.source.table
      ),
      timestamp = getTimestamp(record.payload.source.ts_ms)
    )
  }

  private def getTimestamp(ts: Long): String = {
    OffsetDateTime.ofInstant(Instant.ofEpochMilli(ts), ZoneId.systemDefault()).toString
  }
}
