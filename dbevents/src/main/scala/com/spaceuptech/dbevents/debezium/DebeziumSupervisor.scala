package com.spaceuptech.dbevents.debezium

import akka.actor.typed.Behavior
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.spaceuptech.dbevents.DatabaseSource

object DebeziumSupervisor {
  def apply(): Behavior[Command] = Behaviors.setup[Command](context => new DebeziumSupervisor(context))

  sealed trait Command

  final case class SetDatabaseSources() extends Command

}

class DebeziumSupervisor(context: ActorContext[DebeziumSupervisor.Command]) extends AbstractBehavior[DebeziumSupervisor.Command](context) {

  import DebeziumSupervisor._

  context.log.info("Staring debezium supervisor")

  // Start a test debezium actor
  context.spawn[Debezium.Command](Debezium(DatabaseSource("project", "alias", "postgres", Map("db" -> "postgres", "schema" -> "test"))), s"engine-project-mysql")

  override def onMessage(msg: Command): Behavior[Command] = {
    msg match {
      case SetDatabaseSources() =>
        this
    }
  }
}
