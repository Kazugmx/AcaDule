package net.kazugmx.acadule

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
import org.jetbrains.exposed.sql.exposedLogger
import java.util.*

fun Application.configureAuth(authService: AuthService) {
    // Please read the jwt property from the config file if you are using EngineMain
    val jwtIssuer = environment.config.property("jwt.issuer").getString()
    val jwtAudience = environment.config.property("jwt.audience").getString()
    val jwtDomain = environment.config.property("jwt.domain").getString()
    val jwtRealm = environment.config.property("jwt.realm").getString()
    val jwtSecret = environment.config.property("jwt.secret").getString()

    authentication {
        jwt("auth-jwt") {
            realm = jwtRealm
            verifier(
                JWT
                    .require(Algorithm.HMAC256(jwtSecret))
                    .withAudience(jwtAudience)
                    .withIssuer(jwtDomain)
                    .build()
            )
            validate { credential ->
                if (credential.payload.audience.contains(jwtAudience)) {
                    if( authService.isUserExists(credential.payload.getClaim("userid").asInt()))
                        JWTPrincipal(credential.payload)
                } else null
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
                val chall = authService.login(login)
                if (chall != null) {
                    val generatedToken =
                        JWT.create()
                            .withAudience(jwtAudience)
                            .withIssuer(jwtIssuer)
                            .withClaim("userid", chall.userID)
                            .withExpiresAt(Date(System.currentTimeMillis() + 60000))
                            .sign(Algorithm.HMAC256(jwtSecret))
                    call.respond(
                        mapOf("status" to "success", "token" to generatedToken)
                    )
                }
                call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
            }
            post("tryToken") {
                val token = call.receiveText()
                log.info("TryToken: $token")
                try {
                    val decoded = JWT.require(Algorithm.HMAC256(jwtSecret)).build()
                        .verify(token.removePrefix("Bearer").trim())
                    call.respond(mapOf("status" to "success", "id" to decoded.getClaim("userid").asInt()))
                    authService.isUserExists(decoded.getClaim("userid").asInt())
                } catch (e: TokenExpiredException) {
                    call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed", "reason" to "Token Expired"))
                    log.info("Token Expired: ${e.expiredOn}")
                }
            }
            authenticate("auth-jwt") {
                post("tryTokenV2") {
                    try {
                        val principal = call.principal<JWTPrincipal>()
                            ?: return@post call.respond(
                                HttpStatusCode.BadRequest,
                                mapOf("status" to "failed", "reason" to "No JWT Principal")
                            )
                        principal.payload.getClaim("userid").asInt().let { id ->
                            if (authService.isUserExists(id))
                                call.respond(
                                    mapOf(
                                        "status" to "success",
                                        "id" to id,
                                        "expiry" to principal.expiresAt
                                    )
                                )
                            else call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
                        }
                    } catch (e: TokenExpiredException) {
                        call.respond(
                            HttpStatusCode.BadRequest,
                            mapOf(
                                "status" to "failed",
                                "reason" to "Token Expired ${e.expiredOn}"
                            )
                        )
                    } catch (e: JWTDecodeException) {
                        call.respond(
                            HttpStatusCode.BadRequest, mapOf(
                                "status" to "failed",
                                "reason" to "Invalid Token ${e.message}"
                            )
                        )
                    } catch (e: Exception) {
                        exposedLogger.error("Failed to process tryTokenV2", e)
                        call.respond(
                            HttpStatusCode.InternalServerError,
                            mapOf("status" to "Internal Server Error")
                        )
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