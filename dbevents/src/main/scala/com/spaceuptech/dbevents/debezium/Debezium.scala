package com.spaceuptech.dbevents.debezium

import java.util.concurrent.Executors

import akka.actor.typed.{Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import com.spaceuptech.dbevents.DatabaseSource
import com.spaceuptech.dbevents.debezium.Debezium.{CheckEngineStatus, Command}

import scala.concurrent.duration._

object Debezium {
  def apply(source: DatabaseSource): Behavior[Command] = {
    Behaviors.withTimers[Command](timers => Behaviors.setup[Command](context => new Debezium(context, timers, source))
    )
  }

  sealed trait Command

  case class CheckEngineStatus() extends Command

}

class Debezium(context: ActorContext[Debezium.Command], timers: TimerScheduler[Debezium.Command], source: DatabaseSource) extends AbstractBehavior[Debezium.Command](context) {
  // Extract name of actor
  private val name = Utils.generateConnectorName(source)
  context.log.info(s"Staring debezium engine $name")

  // Start the debezium engine
  private val executor = Executors.newSingleThreadExecutor
  private var status = Utils.startDebeziumEngine(source, executor)

  // Start task for status check
  timers.startTimerAtFixedRate(name, CheckEngineStatus(), 10.second)

  override def onMessage(msg: Command): Behavior[Command] = {
    // No need to handle any messages
    msg match {
      case Debezium.CheckEngineStatus() =>
        // Try starting the debezium engine again only if it wasn't running already
        if (status.future.isDone || status.future.isCancelled) {
          context.log.info(s"Debezium engine $name is closed. Restaring...")
          status = Utils.startDebeziumEngine(source, executor)
        }
        this
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      // Shutdown the timer
      timers.cancel(name)

      // Shut down the debezium engine
      if (!status.future.isCancelled && !status.future.isDone) {
        context.log.info(s"Closing debezium engine - $name")
        status.engine.close()
        context.log.info(s"Closed debezium engine - $name")
      }

      // Shut down the main executor
      executor.shutdownNow()

      this
  }
}
