package com.spaceuptech.dbevents

import akka.actor.typed.ActorSystem

object EventsApp extends App {
  // Load env variables
  val gatewayUrl: String = scala.util.Properties.envOrElse("GATEWAY_URL", "gateway.space-cloud.svc.cluster.local:4122")
  val adminSecret: String = scala.util.Properties.envOrElse("SC_ADMIN_SECRET", "some-secret")

  // Create the main actor system
  val a = ActorSystem[Nothing](EventsSupervisor(gatewayUrl, adminSecret), "db-events")
}
