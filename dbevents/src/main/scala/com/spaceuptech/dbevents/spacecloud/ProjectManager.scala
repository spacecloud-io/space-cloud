package com.spaceuptech.dbevents.spacecloud

import akka.actor.typed.scaladsl.{AbstractBehavior, ActorContext, Behaviors, TimerScheduler}
import akka.actor.typed._
import com.spaceuptech.dbevents.Global
import com.spaceuptech.dbevents.database.Database

import scala.concurrent.duration._
import scala.concurrent.{ExecutionContextExecutor, Future}
import scala.util._

object ProjectManager {
  val fetchEventingConfigKey: String = "fetch-eventing-config"
  val fetchDatabasesKey: String = "fetch-databases"

  def apply(projectId: String): Behavior[Command] =
    Behaviors.withTimers(timers => Behaviors.setup(context => new ProjectManager(context, timers, projectId)))

  sealed trait Command

  case class FetchEventingConfig() extends Command

  case class FetchDatabaseConfig() extends Command

  case class Stop() extends Command
}

class ProjectManager(context: ActorContext[ProjectManager.Command], timers: TimerScheduler[ProjectManager.Command], projectId: String) extends AbstractBehavior(context) {

  import ProjectManager._

  // Member variables
  var databaseToActor: Map[String, ActorRef[Database.Command]] = Map.empty

  // Start the timer
  timers.startTimerAtFixedRate(fetchEventingConfigKey, FetchEventingConfig(), 1.minute)

  override def onMessage(msg: Command): Behavior[Command] = {
    msg match {
      case FetchEventingConfig() =>
        fetchEventingConfig()
        this

      case FetchDatabaseConfig() =>
        fetchDatabaseConfig()
        this


      case Stop() => Behaviors.stopped
    }
  }

  private def fetchDatabaseConfig(): Unit = {
    implicit val system: ActorSystem[Nothing] = context.system
    implicit val executionContext: ExecutionContextExecutor = system.executionContext

    val response: Future[Response[DatabaseConfig]] = fetchSpaceCloudResource(s"http://${Global.gatewayUrl}/v1/config/projects/$projectId/database/config")
    response.onComplete {
      case Success(value) => processDatabaseConfig(value.result)
      case Failure(ex) => context.log.error(s"Unable to fetch database config for project ($projectId)", ex)
    }
  }

  private def processDatabaseConfig(dbs: Array[DatabaseConfig]): Unit = {
    // Filter all disabled databases
    val filteredDbs: Array[DatabaseConfig] = dbs.filter(db => db.enabled)

    // Create actor for new projects
    for (db <- filteredDbs) {
      if (!databaseToActor.contains(db.dbAlias)) {
        val actor = context.spawn(Database.createActor(projectId, db), s"db-${db.dbAlias}")
        databaseToActor += db.dbAlias -> actor
      }
    }

    databaseToActor = databaseToActor.filter(elem => removeDatabaseIfInactive(dbs, elem._1, elem._2))
  }

  private def removeDatabaseIfInactive(dbs: Array[DatabaseConfig], dbAlias: String, actor: ActorRef[Database.Command]): Boolean = {
    if (!dbs.exists(db => db.dbAlias == dbAlias)) {
      actor ! Database.Stop()
      return false
    }
    true
  }

  private def fetchEventingConfig(): Unit = {
    implicit val system: ActorSystem[Nothing] = context.system
    implicit val executionContext: ExecutionContextExecutor = system.executionContext

    val response: Future[Response[EventingConfig]] = fetchSpaceCloudResource(s"http://${Global.gatewayUrl}/v1/config/projects/$projectId/eventing/config")
    response.onComplete {
      case Success(value) => processEventingConfig(value.result(0))
      case Failure(ex) => context.log.error(s"Unable to fetch eventing config for project ($projectId)", ex)
    }
  }

  private def processEventingConfig(config: EventingConfig): Unit = {
    // Stop and remove all children if eventing is disabled
    if (!config.enabled) {
      timers.cancel(fetchDatabasesKey)
      removeAllChildren()
      return
    }

    // Start the timer if its isn't active already
    if (timers.isTimerActive(fetchDatabasesKey)) {
      timers.startTimerAtFixedRate(fetchDatabasesKey, FetchDatabaseConfig(), 1.minute)
      context.self ! FetchDatabaseConfig()
    }
  }

  private def removeAllChildren(): Unit = {
    for ((_, actor) <- databaseToActor) {
      actor ! Database.Stop()
    }
    databaseToActor = Map.empty
  }

  override def onSignal: PartialFunction[Signal, Behavior[Command]] = {
    case PostStop =>
      timers.cancelAll()
      removeAllChildren()
      this
  }
}
