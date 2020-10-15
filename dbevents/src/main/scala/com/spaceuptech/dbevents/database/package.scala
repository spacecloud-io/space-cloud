package com.spaceuptech.dbevents

import java.util.concurrent.Future

import io.debezium.engine.{ChangeEvent, DebeziumEngine}

package object database {

  case class ChangeRecord(payload: ChangeRecordPayload, project: String, dbAlias: String, dbType: String)

  case class ChangeRecordPayload(op: Option[String], before: Option[Map[String, Any]], after: Option[Map[String, Any]], source: ChangeRecordPayloadSource)

  case class ChangeRecordPayloadSource(name: String, ts_ms: Long, table: String)

  case class DebeziumStatus(error: String, future: Future[_], engine: DebeziumEngine[ChangeEvent[String, String]])
}
