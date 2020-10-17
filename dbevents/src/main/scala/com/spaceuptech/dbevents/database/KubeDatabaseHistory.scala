package com.spaceuptech.dbevents.database

import java.util
import java.util.function.Consumer

import io.debezium.config.Configuration
import io.debezium.document.{DocumentReader, DocumentWriter}
import io.debezium.relational.history.{AbstractDatabaseHistory, DatabaseHistoryListener, HistoryRecord, HistoryRecordComparator, TableChanges}
import io.debezium.util.FunctionalReadWriteLock
import io.kubernetes.client.openapi.{ApiClient, ApiException}
import io.kubernetes.client.openapi.apis.CoreV1Api
import io.kubernetes.client.openapi.models.V1ConfigMapBuilder
import io.kubernetes.client.util.ClientBuilder

class KubeDatabaseHistory extends AbstractDatabaseHistory {
  var name: String = ""

  // Helper utilities
  val writer: DocumentWriter = DocumentWriter.defaultWriter()
  val reader: DocumentReader = DocumentReader.defaultReader()
  val lock: FunctionalReadWriteLock = FunctionalReadWriteLock.reentrant()

  // Create kubernetes client
  val client: ApiClient = ClientBuilder.cluster().build()
  io.kubernetes.client.openapi.Configuration.setDefaultApiClient(client)

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
        .withData(new util.HashMap[String, String]())
        .build()
      api.createNamespacedConfigMap("space-cloud", configMap, null, null, null)
    } catch {
      case ex: ApiException => println("Unable to create config map for offset storage", ex.getMessage)
    }
  }

  override def configure(config: Configuration, comparator: HistoryRecordComparator, listener: DatabaseHistoryListener, useCatalogBeforeSchema: Boolean): Unit = {
    super.configure(config, comparator, listener, useCatalogBeforeSchema)
    name = config.getString("name")
    name = name.replaceAll("_","-").toLowerCase + "-history"
    println()
    println("************************************")
    println("Database history store name:", name)
    println("************************************")
    println()
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
        data.put("data", data.get("data") + writer.write(record.document()) + "\n")
        configMap.setData(data)
        api.replaceNamespacedConfigMap(name, "space-cloud", configMap, null, null, null)
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
        val rows = configMap.getData.get("data").split("\n")
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
      val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
      if (configMap == null) {
        return false
      }

      !configMap.getData.isEmpty
    } catch {
      case _: Throwable => false
    }
  }

  override def storageExists(): Boolean = {
    exists()
  }
}
