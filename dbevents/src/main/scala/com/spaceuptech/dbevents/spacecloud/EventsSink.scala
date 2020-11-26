package com.spaceuptech.dbevents.spacecloud

import java.time.{Instant, OffsetDateTime, ZoneId}

import akka.actor.typed.{ActorSystem, Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import com.spaceuptech.dbevents.database.Database.ChangeRecord

import scala.concurrent.ExecutionContextExecutor
import scala.concurrent.duration.DurationInt
import scala.util.{Failure, Success}

object EventsSink {

  def apply(projectId: String): Behavior[Command] = Behaviors.withTimers(timers => Behaviors.setup(context => new EventsSink(context, timers, projectId)))

  sealed trait Command
  case class EmitEvent(record: ChangeRecord) extends Command
  case class EmptyBuffer() extends Command
  case class Stop() extends Command
}

class EventsSink(context: ActorContext[EventsSink.Command], timers: TimerScheduler[EventsSink.Command], projectId: String) extends AbstractBehavior[EventsSink.Command](context) {
  implicit val system: ActorSystem[Nothing] = context.system
  implicit val executionContext: ExecutionContextExecutor = system.executionContext

  timers.startTimerAtFixedRate("send-events", EventsSink.EmptyBuffer(), 250.milliseconds)

  // Create a buffer to hold messages
  var buffer: Array[QueueEvent] = Array.empty

  override def onMessage(msg: EventsSink.Command): Behavior[EventsSink.Command] = {
    msg match {
      case EventsSink.EmitEvent(record) =>
        buffer = buffer :+ prepareQueueRequest(record)

          // Queue events in gateway if the buffer length exceeds 100
        if (buffer.length > 100) {
          val c = buffer
          queueEvent(projectId, c).onComplete {
            case Success(_) => println(s"Events logged successfully - ${c.length}")
            case Failure(exception) => println(s"Unable to log event - ${exception.getMessage}")
          }

          // Empty the buffer
          buffer = Array.empty
        }
        this

      case EventsSink.EmptyBuffer() =>
        // Queue events only if length of buffer exceeds 100
        if (buffer.length > 0) {
          val c = buffer
          queueEvent(projectId, c).onComplete {
            case Success(_) => println(s"Events logged successfully - ${c.length}")
            case Failure(exception) => println(s"Unable to log event - ${exception.getMessage}")
          }

          // Empty the buffer
          buffer = Array.empty
        }
        this

      case EventsSink.Stop() =>
        println(s"Got close command for project event sink - $projectId")
        Behaviors.stopped
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[EventsSink.Command]] = {
    case PostStop =>
      timers.cancelAll()
      buffer = Array.empty
      println(s"Closed project event sink - '$projectId'")
      this
  }

  private def prepareQueueRequest(record: ChangeRecord): QueueEvent = {
    QueueEvent(
      `type` = record.payload.op match {
        case "c" | "r" => "DB_INSERT"
        case "u" => "DB_UPDATE"
        case "d" => "DB_DELETE"
        case _ => s"UNKNOWN_OP_${record.payload.op}"
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
