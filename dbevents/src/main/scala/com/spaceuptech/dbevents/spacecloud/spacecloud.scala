package com.spaceuptech.dbevents

import akka.actor.ClassicActorSystemProvider
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.headers.{Authorization, OAuth2BearerToken}
import akka.http.scaladsl.model.{HttpHeader, HttpMethods, HttpRequest, HttpResponse}
import akka.http.scaladsl.unmarshalling.Unmarshal
import org.json4s.DefaultFormats
import org.json4s.jackson.JsonMethods.parse

import scala.concurrent.{ExecutionContext, Future}

package object spacecloud {

  def fetchSpaceCloudResource[T](url: String)(implicit system: ClassicActorSystemProvider, executor: ExecutionContext, m: Manifest[Response[T]]): Future[Response[T]] = {
    val request = HttpRequest(
      method = HttpMethods.GET,
      uri = url,
      headers = Seq[HttpHeader] {
        Authorization(OAuth2BearerToken(Global.createAdminToken()))
      }
    )

    // Fire the request
    val jsonFuture: Future[String] = Http().singleRequest(request).flatMap {
      case HttpResponse(code, _, entity, _) => {
        val status = code.intValue()
        if (status != 200) {
          Future.failed(new Exception(s"Invalid status code received ($status)"))
        }

        Unmarshal(entity).to[String]
      }
    }

    jsonFuture.flatMap {
      json =>
        implicit val formats: DefaultFormats.type = org.json4s.DefaultFormats
        val res = parse(json).extract[Response[T]]
        if (res.error != "") Future.failed(new Exception(res.error))
        Future {
          res
        }
    }
  }

  case class Response[T](result: Array[T], error: String)

  case class Project(id: String, secret: String)

  case class EventingConfig(enabled: Boolean)

  case class DatabaseConfig(dbAlias: String, `type`: String, name: String, conn: String, enabled: Boolean)

  case class Secret(id: String, data: Map[String, String])
}
