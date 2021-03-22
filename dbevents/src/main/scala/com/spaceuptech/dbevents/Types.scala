package com.spaceuptech.dbevents

case class DatabaseSource(project: String, dbAlias: String, dbType: String, config: Map[String, String])
case class DatabaseSources()
