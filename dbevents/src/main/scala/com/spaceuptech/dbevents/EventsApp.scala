package com.spaceuptech.dbevents

import java.io.File

import akka.actor.typed.ActorSystem
import com.typesafe.config.{Config, ConfigFactory}

object EventsApp extends App {
  // Load env variables
  Global.gatewayUrl = scala.util.Properties.envOrElse("GATEWAY_URL", "gateway.space-cloud.svc.cluster.local:4122")
  Global.secret = scala.util.Properties.envOrElse("SC_ADMIN_SECRET", "some-secret")
  Global.storageType = scala.util.Properties.envOrElse("STORAGE_TYPE", "local")

  val conf = ConfigFactory.load(ConfigFactory.parseFile(new File("/config/application.conf")))
  // Create the main actor system
  val a = ActorSystem[Nothing](EventsSupervisor(), "db-events", conf)
}
