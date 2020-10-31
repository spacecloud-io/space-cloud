package com.spaceuptech.dbevents.pubsub

import java.nio.charset.StandardCharsets

import akka.actor.typed.{ActorSystem, Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.rabbitmq.client.ConnectionFactory
import com.spaceuptech.dbevents.EventsSupervisor
import com.spaceuptech.dbevents.database.Database.ChangeRecord
import com.spaceuptech.dbevents.database.{getConnString}
import org.json4s._
import org.json4s.JsonDSL._
import org.json4s.jackson.JsonMethods._

import scala.concurrent.ExecutionContextExecutor
import scala.util.{Failure, Success}

object RabbitMQ {
  def apply(projectId: String): Behavior[Command] = Behaviors.setup[Command](context => new RabbitMQ(context, projectId))

  sealed trait Command
  case class EmitEvent(record: ChangeRecord) extends Command
  case class UpdateConfig(conn: String) extends Command
  case class Stop() extends Command
}

class RabbitMQ(context: ActorContext[RabbitMQ.Command], projectId: String) extends AbstractBehavior[RabbitMQ.Command](context) {
  import RabbitMQ._

  // Member variables
  private var buffer: Array[ChangeRecord] = Array.empty
  private var connString: String = ""
  private var rabbit: Option[RabbitMQConnection] = None

  override def onMessage(msg: Command): Behavior[Command] = {
    msg match {
      case EmitEvent(record) =>
        send(record)
        this

      case UpdateConfig(conn) =>
        // Implicits we will need
        implicit val system: ActorSystem[Nothing] = context.system
        implicit val executionContext: ExecutionContextExecutor = system.executionContext

        getConnString(projectId, conn) onComplete {
          case Success(c) =>
            // Check if connString has changed
            if (connString == c) return this

            // Connect to rabbitmq
            connect(c)

          case Failure(ex) =>
            context.log.error(s"Unable to get connection string for rabbitmq ($projectId) - ${ex.getMessage}")
        }
        this

      case Stop() => Behaviors.stopped
    }
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      // Close rabbitmq client
      closeClient()
      this
  }

  private def connect(c: String): Unit = {
    connString = c

    // Close the previous client
    closeClient()

    // Prepare connection string
    val factory = new ConnectionFactory()
    factory.setUri(connString)

    // Create connection and channel
    val conn = factory.newConnection()
    val channel = conn.createChannel()

    // Store client for future use
    rabbit = Some(RabbitMQConnection(conn, channel))

    // Empty client
    for (ev <- buffer) {
      send(ev)
    }

    buffer = Array.empty
  }

  private def send(ev: ChangeRecord): Unit = {
    rabbit match {
      case Some(client) =>
        val exchangeName = getExchangeName(ev)

        try {
          implicit lazy val serializerFormats: DefaultFormats.type = org.json4s.DefaultFormats
          val json = compact(render(Extraction.decompose(ev)))

          client.ch.exchangeDeclare(exchangeName, "fanout", false, true, null)
          client.ch.basicPublish(exchangeName, "", null, json.getBytes(StandardCharsets.UTF_8))
        } catch {
          case ex: Throwable => context.log.error(s"Unable to publish message ($exchangeName) - ${ex.getMessage}")
        }
      case None =>
        buffer :+= ev
    }
  }

  private def closeClient(): Unit = {
    rabbit match {
      case Some(value) =>
        // Close the client
        value.ch.close()
        value.conn.close()

        // Reset the client
        rabbit = None

      case None =>
        // No need to do anything here
    }
  }
}
