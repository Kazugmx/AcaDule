package net.kazugmx.acadule.components

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import com.auth0.jwt.exceptions.JWTDecodeException
import com.auth0.jwt.exceptions.TokenExpiredException
import io.ktor.http.HttpStatusCode
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import net.kazugmx.acadule.schemas.AuthService
import net.kazugmx.acadule.schemas.LoginReq
import net.kazugmx.acadule.schemas.UserCreateReq
import java.util.*


suspend inline fun ApplicationCall.safeJwt(block: () -> Unit) {
    try {
        block()
    } catch (e: TokenExpiredException) {
        respond(
            HttpStatusCode.BadRequest,
            mapOf(
                "status" to "failed",
                "reason" to "Token Expired",
                "expiredOn" to e.expiredOn.toString()
            )
        )
        application.log.info("Token Expired: ${e.expiredOn}")
    } catch (e: JWTDecodeException) {
        respond(
            HttpStatusCode.BadRequest,
            mapOf(
                "status" to "failed",
                "reason" to "Invalid Token",
                "detail" to e.message
            )
        )
    } catch (e: Exception) {
        application.log.error("JWT processing failed", e)
        respond(
            HttpStatusCode.InternalServerError,
            mapOf(
                "status" to "failed",
                "reason" to "Internal Server Error"
            )
        )
    }
}

fun Application.configureAuth(authService: AuthService) {
    // Please read the jwt property from the config file if you are using EngineMain
    @Suppress("unused")
    val jwtDomain = environment.config.property("jwt.domain").getString()
    val jwtIssuer = environment.config.property("jwt.issuer").getString()
    val jwtAudience = environment.config.property("jwt.audience").getString()
    val jwtRealm = environment.config.property("jwt.realm").getString()
    val jwtSecret = environment.config.property("jwt.secret").getString()

    authentication {
        jwt("auth-jwt") {
            realm = jwtRealm
            verifier(
                JWT
                    .require(Algorithm.HMAC256(jwtSecret))
                    .withAudience(jwtAudience)
                    .withIssuer(jwtIssuer)
                    .build()
            )
            validate { credential ->
                if (credential.payload.audience.contains(jwtAudience)) {
                    if (authService.isUserExists(credential.payload.getClaim("userid").asInt()))
                        JWTPrincipal(credential.payload) else null
                } else null
            }
            challenge { _, _ ->
                call.respond(HttpStatusCode.Unauthorized, "Token is not valid or has expired")
            }
        }
    }

    routing {
        route("/auth") {
            post("createUser") {
                val user = call.receive<UserCreateReq>()
                val result = authService.createUser(user)
                if (result.status)
                    call.respond(result)
                else
                    call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
            }
            post("login") {
                val login = call.receive<LoginReq>()
                val challenge = authService.login(login)
                if (challenge != null) {
                    val generatedToken =
                        JWT.create()
                            .withAudience(jwtAudience)
                            .withIssuer(jwtIssuer)
                            .withClaim("userid", challenge.userID)
                            .withExpiresAt(
                                Date(
                                    System.currentTimeMillis()
                                            + (1000 * 60 * 60 * 24 * 14)
                                )
                            )
                            .sign(Algorithm.HMAC256(jwtSecret))
                    call.respond(
                        mapOf("status" to "success", "token" to generatedToken)
                    )
                }
                call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
            }
            post("tryToken") {
                call.safeJwt {
                    val token = call.receiveText()
                    val decoded = JWT.require(Algorithm.HMAC256(jwtSecret)).build()
                        .verify(token.removePrefix("Bearer").trim())
                    if (authService.isUserExists(decoded.getClaim("userid").asInt()))
                        call.respond(
                            mapOf(
                                "status" to "success",
                                "id" to decoded.getClaim("userid").toString(),
                                "expiresOn" to decoded.expiresAt.toString()
                            )
                        )
                }
            }
            authenticate("auth-jwt") {
                post("tryTokenV2") {
                    call.safeJwt {
                        val principal = call.principal<JWTPrincipal>()
                            ?: return@post call.respond(
                                HttpStatusCode.BadRequest,
                                mapOf("status" to "failed", "reason" to "No JWT Principal")
                            )
                        principal.payload.getClaim("userid").asInt()?.let { id ->
                            if (authService.isUserExists(id))
                                call.respond(
                                    mapOf(
                                        "status" to "success",
                                        "id" to id.toString(),
                                        "expiry" to principal.expiresAt.toString()
                                    )
                                )
                            else call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
                        }
                    }
                }
                get("users") {
                    call.respond(
                        authService.getUsers()
                    )
                }
            }
        }
    }
}