package com.spaceuptech.spacecloud

import akka.actor.ClassicActorSystemProvider
import akka.actor.typed.ActorSystem
import akka.http.scaladsl.Http
import akka.http.scaladsl.model.{ContentTypes, HttpEntity, HttpRequest}

import scala.util.Success

object Utils {
  def getProjectDetails(gatewayUrl: String, project: String)(implicit system: ClassicActorSystemProvider): Unit = {
    val request = HttpRequest(uri = s"http://$gatewayUrl/v1/config/projects/$project")
    val response = Http().singleRequest(request)
    response {
      case Success(res) =>
    }
  }

  def prepareSpaceCloudGraphQLUrl(gatewayUrl: String, project: String): String = {
    s"http://$gatewayUrl/v1/api/$project/graphql"
  }
}
