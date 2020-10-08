package com.spaceuptech.dbevents

import akka.actor.typed.{Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.spaceuptech.dbevents.debezium.DebeziumSupervisor
import com.spaceuptech.spacecloud.SpaceCloudClient

// We are creating an object to simply creation of the actor
object EventsSupervisor {
  def apply(gatewayUrl: String, adminSecret: String): Behavior[Nothing] = Behaviors.setup[Nothing](context => new EventsSupervisor(context, gatewayUrl, adminSecret))
}

class EventsSupervisor(context: ActorContext[Nothing], gatewayUrl: String, adminSecret: String) extends AbstractBehavior[Nothing](context) {
  context.log.info("DB events source app started")

  // Start the debezium supervisor
  private val debezium = context.spawn(DebeziumSupervisor(), "debezium")

  // Start the space cloud client
  private val spaceCloudClient = context.spawn(SpaceCloudClient(gatewayUrl, adminSecret), "space-cloud")


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
