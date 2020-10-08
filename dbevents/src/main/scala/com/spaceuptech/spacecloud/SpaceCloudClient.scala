package com.spaceuptech.spacecloud

import akka.actor.typed.{ActorRef, ActorSystem, Behavior}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import akka.http.scaladsl.Http
import akka.http.scaladsl.model._
import org.json4s._
import org.json4s.jackson.Serialization

import scala.concurrent.{ExecutionContextExecutor, Future}
import scala.util.{Failure, Success}

object SpaceCloudClient {
  def apply(gatewayUrl: String, adminSecret: String): Behavior[Command] = Behaviors.setup[Command](context => new SpaceCloudClient(context, gatewayUrl, adminSecret))

  sealed trait Command
  case class GraphQLRequest(id: String, project: String, query: String, variables: Map[String, AnyVal], replyTo: ActorRef[GraphQLResponse]) extends Command
  case class GraphQLResponse(id: String, ack: Boolean, response: Any)
}

class SpaceCloudClient(context: ActorContext[SpaceCloudClient.Command], gatewayUrl: String, adminSecret: String) extends AbstractBehavior[SpaceCloudClient.Command](context) {
  import SpaceCloudClient._

  override def onMessage(msg: Command): Behavior[Command] = {
    msg match {
      case GraphQLRequest(id, project, query, variables, replyTo) =>
        // Prepare json string
        implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
        val jsonString = Serialization.write(Map[String, Any]("query" -> query, "variables" -> variables))

        // Prepare http request
        val request = HttpRequest(
          method = HttpMethods.POST,
          uri = Utils.prepareSpaceCloudGraphQLUrl(gatewayUrl, project),
          entity = HttpEntity(ContentTypes.`application/json`, jsonString)
        )

        // Fire the request
        implicit val system: ActorSystem[Nothing] = context.system
        implicit val executionContext: ExecutionContextExecutor = context.system.executionContext
        val response: Future[HttpResponse] = Http().singleRequest(request)

        response.onComplete {
          case Success(res) => replyTo ! GraphQLResponse(id, ack = true, res)
          case Failure(_) =>  replyTo ! GraphQLResponse(id, ack = false, null)
        }

        this
    }
  }


}
