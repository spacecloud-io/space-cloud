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
import org.apache.kafka.connect.runtime.WorkerConfig
import org.apache.kafka.connect.storage.OffsetBackingStore
import org.apache.kafka.connect.util.Callback

import scala.jdk.CollectionConverters._


class KubeOffsetBackingStore extends OffsetBackingStore {
  // Variable to store identity of caller
  var name: String = ""

  var executor: ExecutorService = _

  // Create kubernetes client
  val client: ApiClient = ClientBuilder.cluster().build()
  Configuration.setDefaultApiClient(client)

  override def start(): Unit = {
    executor = Executors.newSingleThreadExecutor

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

  override def stop(): Unit = {
    // Stop the executor first
    if (executor != null) {
      executor.shutdown()
      try {
        executor.awaitTermination(30, TimeUnit.SECONDS)
      } catch {
        case ex: Throwable => println("unable to stop the executor:", ex.getMessage)
      }
    }

    // Create a v1 api client
    val api = new CoreV1Api()

    // Go ahead and delete all config maps associated with this connector
    try {
      api.deleteNamespacedConfigMap(name, "space-cloud", null, null, null, null, null, new V1DeleteOptions())
    } catch {
      case ex: Throwable => println("Unable tot delete config maps:", ex.getMessage)
    }
  }

  override def get(keys: util.Collection[ByteBuffer]): Future[util.Map[ByteBuffer, ByteBuffer]] = {
    executor.submit(() => {
      // Make a result map
      val result: util.Map[ByteBuffer, ByteBuffer] = new util.HashMap()

      // Create a v1 api client
      val api = new CoreV1Api()

      // Get the config map
      val configMap = api.readNamespacedConfigMap(name, "space-cloud", null,null, null)

      // Iterate over the keys
      val itr = keys.iterator()
      while(itr.hasNext) {
        val key = itr.next()
        val value = configMap.getData.get(Base64.getEncoder.encodeToString(key.array()))
        result.put(key, StandardCharsets.UTF_8.encode(value))
      }

      result
    })
  }

  override def set(values: util.Map[ByteBuffer, ByteBuffer], callback: Callback[Void]): Future[Void] = {
    executor.submit(() => {
      try {
        // Create a v1 api client
        val api = new CoreV1Api()

        // Get the config map
        val configMap = api.readNamespacedConfigMap(name, "space-cloud", null, null, null)
        val currentValues = configMap.getData

        // Store the values in the config map
        val map = values.asScala
        for ((k, v) <- map) {
          currentValues.put(Base64.getEncoder.encodeToString(k.array()), Base64.getEncoder.encodeToString(v.array()))
        }
        configMap.setData(currentValues)

        // Update the config map
        api.replaceNamespacedConfigMap(name, "space-cloud", configMap, null, null, null)

        if (callback != null) callback.onCompletion(null, null)
      } catch {
        case ex: Throwable => if (callback != null) callback.onCompletion(ex, null)
      }

      null
    })
  }

  override def configure(config: WorkerConfig): Unit = {
    name = config.getString("name")
    name = name.replaceAll("_","-").toLowerCase
    println()
    println("************************************")
    println("Offset backing store name:", name)
    println("************************************")
    println()
  }
}
