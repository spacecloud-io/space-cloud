package com.spaceuptech.dbevents.database

import java.nio.ByteBuffer
import java.nio.charset.StandardCharsets

import com.spaceuptech.dbevents.Global
import org.bson.BsonDocument

class MongoStore {
  private val store = Global.storageType match {
    case "k8s" => Some(new KubeOffsetBackingStore())
    case _ => None
  }


  def setName(name: String): Unit = {
    store match {
      case Some(value) => value.setName(name)
      case None =>
    }
  }

  def start(): Unit = {
    try {
      store match {
        case Some(value) => value.start()
        case None =>
      }
    } catch {
      case ex: Throwable => println(s"Unable to start mongo store resume token - ${ex.getMessage}")
    }
  }

  def stop(): Unit = {
    try {
      store match {
        case Some(value) => value.stop()
        case None =>
      }
    } catch {
      case ex: Throwable => println(s"Unable to stop mongo store resume token - ${ex.getMessage}")
    }
  }

  def get(): Option[BsonDocument] = {
    store match {
      case Some(value) =>
        try {
          val key = StandardCharsets.UTF_8.encode("resume-token")
          val resumeToken = value.get(java.util.Arrays.asList(key)).get().get(key)

          // Return the resume token if it isn't null
          if (resumeToken == null) {
            return Some(mongoByteBufferToBsonDocument(resumeToken))
          }
        } catch {
          case ex: Throwable => println(s"Unable to get mongo store resume token - ${ex.getMessage}")
        }

        None

      case None => None
    }
  }

  def set(resumeToken: BsonDocument): Unit = {
    store match {
      case Some(k8sStore) =>
        try {
          val key = StandardCharsets.UTF_8.encode("resume-token")
          val value = StandardCharsets.UTF_8.encode(resumeToken.toJson)
          k8sStore.set(java.util.Collections.singletonMap(key, value), null)
        } catch {
          case ex: Throwable => println(s"Unable to set mongo store resume token - ${ex.getMessage}")
        }

      case None =>
    }
  }

  private def mongoByteBufferToBsonDocument(data: ByteBuffer): BsonDocument = {
    BsonDocument.parse(new String(data.array(), "UTF-8"))
  }
}
