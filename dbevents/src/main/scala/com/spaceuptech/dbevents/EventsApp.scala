package com.spaceuptech.dbevents

import akka.actor.typed.ActorSystem

object EventsApp extends App {
  // Create the main actor system
 val a = ActorSystem[Nothing](EventsSupervisor(), "db-events")
}
