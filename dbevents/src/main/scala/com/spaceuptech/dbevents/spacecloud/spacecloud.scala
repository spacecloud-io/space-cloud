package com.spaceuptech.dbevents

import akka.actor.ClassicActorSystemProvider
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.headers.{Authorization, OAuth2BearerToken}
import akka.http.scaladsl.model.{ContentTypes, HttpEntity, HttpHeader, HttpMethods, HttpRequest, HttpResponse, StatusCodes}
import akka.http.scaladsl.unmarshalling.Unmarshal
import org.json4s.{DefaultFormats, Extraction}
import org.json4s.jackson.JsonMethods._

import scala.concurrent.{ExecutionContext, Future}

package object spacecloud {

  def queueEvent(project: String, event: QueueEvent)(implicit system: ClassicActorSystemProvider, executor: ExecutionContext): Future[Unit] = {
    implicit lazy val serializerFormats: DefaultFormats.type = org.json4s.DefaultFormats
    val request = HttpRequest(
      method = HttpMethods.POST,
      uri = s"http://${Global.gatewayUrl}/v1/api/$project/eventing/admin-queue",
      headers = Seq[HttpHeader] {
        Authorization(OAuth2BearerToken(Global.createAdminToken()))
      },
      entity = HttpEntity(
        ContentTypes.`application/json`,
        compact(render(Extraction.decompose(event)))
      )
    )

    Http().singleRequest(request).flatMap {
      case resp @ HttpResponse(StatusCodes.OK, _, _, _) =>
        resp.discardEntityBytes()
        Future{}

      case HttpResponse(code, _, _, _) =>
        Future.failed(new Exception(s"Invalid status code received (${code.intValue()})"))
    }
  }

  def fetchSpaceCloudResource[T: Manifest](url: String)(implicit system: ClassicActorSystemProvider, executor: ExecutionContext): Future[T] = {
    val request = HttpRequest(
      method = HttpMethods.GET,
      uri = url,
      headers = Seq[HttpHeader] {
        Authorization(OAuth2BearerToken(Global.createAdminToken()))
      }
    )

    // Fire the request
    val jsonFuture: Future[String] = Http().singleRequest(request).flatMap {
      case HttpResponse(code, _, entity, _) =>
        val status = code.intValue()
        if (status != 200) {
          Future.failed(new Exception(s"Invalid status code received ($status)"))
        }

        Unmarshal(entity).to[String]
    }

    jsonFuture.flatMap {
      json =>
        implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
        val res = parse(json).extract[T]
        Future {
          res
        }
    }
  }


  case class ProjectResponse(result: Array[Project], error: Option[String])
  case class Project(id: String)

  case class EventingConfigResponse(result: Array[EventingConfig], error: Option[String])
  case class EventingConfig(enabled: Boolean, dbAlias: String)

  case class DatabaseConfigResponse(result: Array[DatabaseConfig], error: Option[String])
  case class DatabaseConfig(dbAlias: String, `type`: String, name: String, conn: String, enabled: Boolean)

  case class SecretResponse(result: Array[Secret], error: Option[String])
  case class Secret(id: String, data: Map[String, String])

  case class QueueEvent(`type`: String, timestamp: String, payload: DatabaseEvent, options: DatabaseEventOptions)

  case class DatabaseEvent(db: String, col: String, doc: Option[Map[String, Any]], find: Option[Map[String, Any]], before: Option[Map[String, Any]])

  case class DatabaseEventOptions(db: String, col: String)
}
