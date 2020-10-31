package com.spaceuptech.dbevents

import com.rabbitmq.client.{Channel, Connection}
import com.spaceuptech.dbevents.database.Database.ChangeRecord

package object pubsub {
  case class RabbitMQConnection(conn: Connection, ch: Channel)

  def getExchangeName(event: ChangeRecord): String = {
    s"${event.project}___${event.dbAlias}___${event.payload.source.table}___${event.payload.op}"
  }
}
