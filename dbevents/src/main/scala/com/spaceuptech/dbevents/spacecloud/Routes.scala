package com.spaceuptech.dbevents.spacecloud

import akka.actor.typed.{ActorRef, ActorSystem}
import akka.http.scaladsl.model._
import akka.http.scaladsl.server.Directives._
import akka.http.scaladsl.server.Route
import akka.util.Timeout

import scala.concurrent.duration.DurationInt

class Routes(projects: ActorRef[ProjectsSupervisor.Command])(implicit system: ActorSystem[_]) {
  implicit  val timeout: Timeout = 3.seconds

  lazy val routes: Route =
    path("fetch-projects") {
      post {
        projects ! ProjectsSupervisor.FetchProjects()
        complete(HttpEntity(ContentTypes.`application/json`, "{}"))
      }
    }
}
