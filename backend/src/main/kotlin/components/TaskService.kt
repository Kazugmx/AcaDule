package net.kazugmx.acadule.components

import io.ktor.http.HttpStatusCode
import io.ktor.server.application.Application
import io.ktor.server.application.log
import io.ktor.server.auth.authenticate
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.principal
import io.ktor.server.request.receive
import io.ktor.server.response.respond
import io.ktor.server.routing.delete
import io.ktor.server.routing.get
import io.ktor.server.routing.patch
import io.ktor.server.routing.post
import io.ktor.server.routing.route
import io.ktor.server.routing.routing
import net.kazugmx.acadule.schemas.AuthService
import net.kazugmx.acadule.schemas.CreateTaskReq
import net.kazugmx.acadule.schemas.IDTaskReq
import net.kazugmx.acadule.schemas.TaskService
import net.kazugmx.acadule.schemas.UpdateTaskReq
import org.jetbrains.exposed.sql.Database
import java.util.*

fun Application.configureTaskService(
    database: Database,
    authService: AuthService,
) {
    // TODO task service
    val taskService = TaskService(database)
    routing {
        authenticate("auth-jwt") {
            route("/task") {
                get {
                    call.safeJwt {
                        val principal =
                            call.principal<JWTPrincipal>() ?: return@get call.respond(
                                HttpStatusCode.BadRequest,
                                mapOf("status" to "failed", "reason" to "No JWT Principal"),
                            )
                        val id = principal.payload.getClaim("userid").asInt()
                        call.respond(
                            HttpStatusCode.OK,
                            taskService.getAllTasksByUser(id),
                        )
                    }
                }
                post {
                    call.safeJwt { principal ->
                        val taskPayload = call.receive<CreateTaskReq>()
                        principal.payload.getClaim("userid").asInt()?.let { id ->
                            if (authService.isUserExists(id)) {
                                taskService
                                    .createTask(
                                        taskPayload,
                                        id,
                                    ).value
                                    .toString()
                                    .let {
                                        call.respond(
                                            HttpStatusCode.OK,
                                            mapOf("status" to "success", "id" to it),
                                        )
                                    }
                            }
                        }
                    }
                }
                patch {
                    call.safeJwt { principal ->
                        val patchingTask = call.receive<UpdateTaskReq>()
                        val userID = principal.payload.getClaim("userid").asInt()
                        val updated =
                            taskService.updateTask(patchingTask, userID)
                                ?: return@patch call.respond(
                                    HttpStatusCode.NotFound,
                                    mapOf("status" to "failed", "reason" to "Task not found"),
                                )
                        call.respond(HttpStatusCode.OK, updated)
                    }
                }

                delete {
                    call.safeJwt { principal ->
                        val deletingTask = call.receive<IDTaskReq>()
                        val userID = principal.payload.getClaim("userid").asInt()
                        if (authService.isUserExists(userID)) {
                            taskService.deleteTask(id = deletingTask.id, userID = userID)
                        }
                    }
                }
                route("/{taskID}") {
                    get {
                        call.safeJwt { principal ->
                            val id = principal.payload.getClaim("userid").asInt()
                            try {
                                val taskId = UUID.fromString(call.parameters["taskID"])
                                log.info("Task ID: $taskId")
                                val task =
                                    taskService.getTaskByID(taskID = taskId.toString(), userID = id) ?: call.respond(
                                        HttpStatusCode.NotFound,
                                        mapOf("status" to "failed", "reason" to "Task not found"),
                                    )
                                call.respond(HttpStatusCode.OK, task)
                            } catch (_: IllegalArgumentException) {
                                call.respond(
                                    HttpStatusCode.BadRequest,
                                    mapOf("status" to "failed", "reason" to "Invalid Task ID"),
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}
