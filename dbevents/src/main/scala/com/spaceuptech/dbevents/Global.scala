package com.spaceuptech.dbevents

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm

object Global {
  var secret: String = ""
  var gatewayUrl: String = ""
  var storageType: String = "local"

  def createAdminToken(): String = {
    val alg = Algorithm.HMAC256(secret)
    JWT.create()
      .withClaim("role", "internal")
      .withClaim("id", "debezium")
      .sign(alg)
  }
}
