package net.kazugmx.acadule

import io.ktor.server.application.*
import io.ktor.server.request.receive
import io.ktor.server.routing.*
import net.kazugmx.acadule.schemas.TaskService
import org.jetbrains.exposed.sql.Database

fun Application.configureDatabases(database: Database) {
    val taskService = TaskService(database)

    routing {
        route("/tasks"){
            post ("createTask"){
            }
        }
    }
}