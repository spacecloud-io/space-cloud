package com.spaceuptech.dbevents

import akka.actor.typed.{Behavior, PostStop, Signal}
import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors}
import com.spaceuptech.dbevents.spacecloud.ProjectsSupervisor

// We are creating an object to simply creation of the actor
object EventsSupervisor {
  def apply(): Behavior[Nothing] = Behaviors.setup[Nothing](context => new EventsSupervisor(context))
}

class EventsSupervisor(context: ActorContext[Nothing]) extends AbstractBehavior[Nothing](context) {
  println("DB events source app started")

  // Start the projects supervisor
  private val projects = context.spawn(ProjectsSupervisor(), "projects")
  projects ! ProjectsSupervisor.FetchProjects()

  // No need to handle any messages
  override def onMessage(msg: Nothing): Behavior[Nothing] = Behaviors.unhandled

  override def onSignal: PartialFunction[Signal, Behavior[Nothing]] = {
    case PostStop =>
      projects ! ProjectsSupervisor.Stop()
      println("DB events source app stopped")
      this
  }
}
