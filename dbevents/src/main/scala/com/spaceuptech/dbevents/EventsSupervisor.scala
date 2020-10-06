package com.spaceuptech.dbevents

import akka.actor.typed.{Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.spaceuptech.dbevents.debezium.DebeziumSupervisor

// We are creating an object to simply creation of the actor
object EventsSupervisor {
  def apply(): Behavior[Nothing] = Behaviors.setup[Nothing](context => new EventsSupervisor(context))
}

class EventsSupervisor(context: ActorContext[Nothing]) extends AbstractBehavior[Nothing](context) {
  context.log.info("DB events source app started")

  // Start the debezium supervisor
  private val debezium = context.spawn(DebeziumSupervisor(), "debezium")

  override def onMessage(msg: Nothing): Behavior[Nothing] = {
    // No need to handle any messages
    Behaviors.unhandled
  }

  override def onSignal: PartialFunction[Signal, Behavior[Nothing]] = {
    case PostStop =>
      context.log.info("DB events source app stopped")
      this
  }
}
