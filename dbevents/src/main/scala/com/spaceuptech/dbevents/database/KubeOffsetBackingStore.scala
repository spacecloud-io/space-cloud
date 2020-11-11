package com.spaceuptech.dbevents.database

import java.nio.ByteBuffer
import java.nio.charset.StandardCharsets
import java.util
import java.util.Base64
import java.util.concurrent.{Executor, ExecutorService, Executors, Future, TimeUnit}

import io.kubernetes.client.openapi.apis.CoreV1Api
import io.kubernetes.client.openapi.models.{V1ConfigMap, V1ConfigMapBuilder, V1DeleteOptions}
import io.kubernetes.client.openapi.{ApiClient, ApiException, Configuration}
import io.kubernetes.client.util.ClientBuilder
import org.apache.commons.codec.binary.Hex
import org.apache.commons.lang3.concurrent.ConcurrentUtils
import org.apache.kafka.connect.runtime.WorkerConfig
import org.apache.kafka.connect.storage.OffsetBackingStore
import org.apache.kafka.connect.util.Callback

import scala.jdk.CollectionConverters._


class KubeOffsetBackingStore extends OffsetBackingStore {
  // Variable to store identity of caller
  var name: String = ""

  // Create kubernetes client
  val client: ApiClient = ClientBuilder.cluster().build()

  override def start(): Unit = {

    // Check if the store has already been configured
    if (name == "") {
      throw new Exception("Call configure before calling start")
    }

    // Create an empty config map if it doesn't already exist
    try {
      println("Creating config map for offset:", name)

      val client: ApiClient = ClientBuilder.cluster().build()

      val api = new CoreV1Api()
      api.setApiClient(client)

      val configMap = new V1ConfigMapBuilder()
        .withNewMetadata().withName(name).withLabels(Map("app" -> "debezium").asJava).endMetadata()
        .withData(Map("test" -> "key").asJava)
        .build()
      api.createNamespacedConfigMap("space-cloud", configMap, null, null, null)
    } catch {
      case ex: Throwable => println("Unable to create config map for offset storage", ex.getMessage)
    }
  }

  override def stop(): Unit = {

    // Create a v1 api client
    val api = new CoreV1Api()
    api.setApiClient(client)

    // Go ahead and delete all config maps associated with this connector
    try {
      println("Deleting config map for offset:", name)
      api.deleteNamespacedConfigMap(name, "space-cloud", null, null, null, null, null, new V1DeleteOptions())
    } catch {
      case ex: Throwable => println("Unable to delete config maps:", ex.getMessage)
    }
  }

  override def get(keys: util.Collection[ByteBuffer]): Future[util.Map[ByteBuffer, ByteBuffer]] = {
    // Make a result map
    val result: util.Map[ByteBuffer, ByteBuffer] = new util.HashMap()

    // Create a v1 api client
    val api = new CoreV1Api()
    api.setApiClient(client)

    // Get the config map
    println("Getting config map for offset:", name)
    val configMap = api.readNamespacedConfigMap(name, "space-cloud", null,null, null)

    // Iterate over the keys
    val itr = keys.iterator()
    while(itr.hasNext) {
      val key = itr.next()
      val value = configMap.getData.get(Hex.encodeHexString(key.array(), true))
      if (value != null) {
        result.put(key, ByteBuffer.wrap(Hex.decodeHex(value)))
      } else {
        result.put(key, null)
      }
    }

    ConcurrentUtils.constantFuture(result)
  }

  override def set(values: util.Map[ByteBuffer, ByteBuffer], callback: Callback[Void]): Future[Void] = {
    try {
      // Create a v1 api client
      val api = new CoreV1Api()
      api.setApiClient(client)

      // Get the config map
      println("Getting config map to set for offset:", name)
      val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
      val currentValues = configMap.getData

      // Store the values in the config map
      val map = values.asScala
      for ((k, v) <- map) {
        currentValues.put(Hex.encodeHexString(k.array(), true), Hex.encodeHexString(v.array(), true))
      }
      configMap.setData(currentValues)

      // Update the config map
      println("Setting config map for offset:", name)
      val configMap2 = new V1ConfigMapBuilder()
        .withNewMetadata().withName(name).withLabels(Map("app" -> "debezium").asJava).endMetadata()
        .withData(currentValues)
        .build()
      api.replaceNamespacedConfigMap(name, "space-cloud", configMap2, null, null, null)

      if (callback != null) callback.onCompletion(null, null)
    } catch {
      case ex: Throwable => if (callback != null) callback.onCompletion(ex, null)
    }

    ConcurrentUtils.constantFuture(null)
  }

  override def configure(config: WorkerConfig): Unit = {
    setName(config.getString("offset.storage.file.filename"))
    println()
    println("************************************")
    println("Offset backing store name:", name)
    println("************************************")
    println()
  }

  def setName(value: String): Unit = {
    name = value.replaceAll("_","-").replaceAll("[.]", "").replaceAll("/", "").toLowerCase
  }
}
