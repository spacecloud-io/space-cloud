package com.spaceuptech.dbevents

import akka.actor.typed.ActorSystem

object EventsApp extends App {
  // Load env variables
  Global.gatewayUrl = scala.util.Properties.envOrElse("GATEWAY_URL", "gateway.space-cloud.svc.cluster.local:4122")
  Global.secret = scala.util.Properties.envOrElse("SC_ADMIN_SECRET", "some-secret")
  Global.storageType = scala.util.Properties.envOrElse("STORAGE_TYPE", "local")

  // Create the main actor system
  val a = ActorSystem[Nothing](EventsSupervisor(), "db-events")
}
