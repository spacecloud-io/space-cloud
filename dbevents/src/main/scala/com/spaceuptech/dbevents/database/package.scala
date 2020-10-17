package com.spaceuptech.dbevents


import akka.actor.ClassicActorSystemProvider
import com.spaceuptech.dbevents.spacecloud.{Secret, fetchSpaceCloudResource}
import io.debezium.engine.{ChangeEvent, DebeziumEngine}

import scala.concurrent.{ExecutionContext, Future}

package object database {

  case class ChangeRecordPayload(op: Option[String], before: Option[Map[String, Any]], after: Option[Map[String, Any]], source: ChangeRecordPayloadSource)
  case class ChangeRecordPayloadSource(name: String, ts_ms: Long, table: String)

  case class DebeziumStatus(error: String, future: java.util.concurrent.Future[_], engine: DebeziumEngine[ChangeEvent[String, String]])

  def getConnString(projectId: String, conn: String)(implicit system: ClassicActorSystemProvider, executor: ExecutionContext): Future[String] = {
    if (!conn.startsWith("secrets")) {
      return Future { conn }
    }

    val secret = conn.split('.')(1)
    fetchSpaceCloudResource[Secret](s"http://${Global.gatewayUrl}/v1/runner/$projectId/secrets?id=$secret").flatMap {
      secretResponse =>
        secretResponse.result(0).data.get("CONN") match {
          case Some(conn) => Future{conn}
          case _ => Future.failed(new Exception("Secret does not have a valid resonse"))
        }
    }
  }
}
