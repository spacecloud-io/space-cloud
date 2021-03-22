package com.spaceuptech.dbevents.database

import java.util
import java.util.function.Consumer

import io.debezium.config.Configuration
import io.debezium.document.{DocumentReader, DocumentWriter}
import io.debezium.relational.history.{AbstractDatabaseHistory, DatabaseHistoryListener, HistoryRecord, HistoryRecordComparator, TableChanges}
import io.debezium.util.FunctionalReadWriteLock
import io.kubernetes.client.openapi.{ApiClient, ApiException}
import io.kubernetes.client.openapi.apis.CoreV1Api
import io.kubernetes.client.openapi.models.{V1ConfigMapBuilder, V1DeleteOptions}
import io.kubernetes.client.util.ClientBuilder

import scala.jdk.CollectionConverters._

class KubeDatabaseHistory extends AbstractDatabaseHistory {
  var name: String = ""

  // Helper utilities
  val writer: DocumentWriter = DocumentWriter.defaultWriter()
  val reader: DocumentReader = DocumentReader.defaultReader()
  val lock: FunctionalReadWriteLock = FunctionalReadWriteLock.reentrant()

  // Create kubernetes client
  val client: ApiClient = ClientBuilder.cluster().build()
  io.kubernetes.client.openapi.Configuration.setDefaultApiClient(client)

  override def stop(): Unit = {
    super.stop()
    // Create a v1 api client
    val api = new CoreV1Api()

    // Go ahead and delete all config maps associated with this connector
    try {
      println("Deleting config map for db history:", name)
      api.deleteNamespacedConfigMap(name, "space-cloud", null, null, null, null, null, new V1DeleteOptions())
    } catch {
      case ex: Throwable => println("Unable to delete config maps for db history:", ex.getMessage)
    }
  }
  override def start(): Unit = {
    super.start()

    // Check if the store has already been configured
    if (name == "") {
      throw new Exception("Call configure before calling start")
    }

    // Create an empty config map if it doesn't already exist
    try {
      val api = new CoreV1Api()
      val configMap = new V1ConfigMapBuilder()
        .withNewMetadata().withName(name).endMetadata()
        .withData(Map("data" -> "").asJava)
        .build()
      println("Create database history config map", name)
      api.createNamespacedConfigMap("space-cloud", configMap, null, null, null)
    } catch {
      case ex: ApiException => println("Unable to create config map for offset storage", ex.getMessage)
    }
  }

  override def configure(config: Configuration, comparator: HistoryRecordComparator, listener: DatabaseHistoryListener, useCatalogBeforeSchema: Boolean): Unit = {
    super.configure(config, comparator, listener, useCatalogBeforeSchema)
    setName(config.getString("database.history.file.filename"))
    println()
    println("************************************")
    println("Database history store name:", name)
    println("************************************")
    println()
  }

  def setName(value: String): Unit = {
    name = value.replaceAll("_","-").replaceAll("[.]", "").replaceAll("/", "").toLowerCase
  }

  override def storeRecord(record: HistoryRecord): Unit = {
    if (record == null) {
      return
    }

    lock.write(new Runnable {
      override def run(): Unit = {
        // Create a v1 api client
        val api = new CoreV1Api()

        // Update the config map
        val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
        val data = configMap.getData
        data.put("data", data.get("data") + writer.write(record.document()) + "----")
        val configMap2 = new V1ConfigMapBuilder()
          .withNewMetadata().withName(name).endMetadata()
          .withData(data)
          .build()
        api.replaceNamespacedConfigMap(name, "space-cloud", configMap2, null, null, null)
      }
    })
  }

  override def recoverRecords(records: Consumer[HistoryRecord]): Unit = {
    lock.write(new Runnable {
      override def run(): Unit = {
        // Create a v1 api client
        val api = new CoreV1Api()

        // Update the config map
        val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
        val rows = configMap.getData.get("data").split("----")
        for (row <- rows) {
          if (row.length != 0) {
            records.accept(new HistoryRecord(reader.read(row)))
          }
        }
      }
    })
  }

  override def exists(): Boolean = {
    // Create a v1 api client
    val api = new CoreV1Api()

    try {
      println("Checking if db history config map exists", name)
      val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
      if (configMap == null || configMap.getData == null) {
        return false
      }

      configMap.getData.get("data").length > 0
    } catch {
      case _: Throwable => false
    }
  }

  override def storageExists(): Boolean = {
    exists()
  }
}
