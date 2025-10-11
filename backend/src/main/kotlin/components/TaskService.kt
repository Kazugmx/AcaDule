package net.kazugmx.acadule.components

import io.ktor.http.HttpStatusCode
import io.ktor.server.application.Application
import io.ktor.server.auth.authenticate
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.principal
import io.ktor.server.response.respond
import io.ktor.server.routing.get
import io.ktor.server.routing.route
import io.ktor.server.routing.routing
import net.kazugmx.acadule.schemas.AuthService
import net.kazugmx.acadule.schemas.TaskService
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.exposedLogger

fun Application.configureTaskService(database: Database, authService: AuthService) {
    //TODO task service
    val taskService = TaskService(database)
    routing {
        authenticate("auth-jwt") {
            route("/tasks") {
                get {
                    try {

                        val principal = call.principal<JWTPrincipal>()
                            ?: return@get call.respond(
                                HttpStatusCode.BadRequest,
                                mapOf("status" to "failed", "reason" to "No JWT Principal")
                            )
                        principal.payload.getClaim("userid").asInt()?.let { id ->
                            if (authService.isUserExists(id)) {
                                call.respond(
                                    HttpStatusCode.OK,
                                    taskService.getAllTasksByUser(id)
                                )
                            } else call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
                        }
                    }
                    catch (e: Exception) {
                        exposedLogger.error("Failed to process request getTasks", e)
                        call.respond(
                            HttpStatusCode.InternalServerError,
                            mapOf("status" to "Internal Server Error")
                        )
                    }
                }
            }
        }
    }
}