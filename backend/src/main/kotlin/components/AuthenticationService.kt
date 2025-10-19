package net.kazugmx.acadule.components

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
import io.ktor.http.HttpStatusCode
import io.ktor.server.application.Application
import io.ktor.server.auth.authenticate
import io.ktor.server.auth.authentication
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.jwt.jwt
import io.ktor.server.request.receive
import io.ktor.server.response.respond
import io.ktor.server.routing.get
import io.ktor.server.routing.post
import io.ktor.server.routing.route
import io.ktor.server.routing.routing
import net.kazugmx.acadule.schemas.AuthService
import net.kazugmx.acadule.schemas.LoginReq
import net.kazugmx.acadule.schemas.UserCreateReq
import java.util.Date

fun Application.configureAuth(authService: AuthService) {
    // Please read the jwt property from the config file if you are using EngineMain
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
                    .build(),
            )
            validate { credential ->
                if (credential.payload.audience.contains(jwtAudience)) {
                    if (authService.isUserExists(credential.payload.getClaim("userid").asInt())) {
                        JWTPrincipal(credential.payload)
                    } else {
                        null
                    }
                } else {
                    null
                }
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
                if (result.status) {
                    call.respond(result)
                } else {
                    call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
                }
            }
            post("login") {
                val login = call.receive<LoginReq>()
                val challenge = authService.login(login)
                if (challenge != null) {
                    val generatedToken =
                        JWT
                            .create()
                            .withAudience(jwtAudience)
                            .withIssuer(jwtIssuer)
                            .withClaim("userid", challenge.userID)
                            .withExpiresAt(
                                Date(
                                    System.currentTimeMillis() +
                                        (1000 * 60 * 60 * 24 * 14),
                                ),
                            ).sign(Algorithm.HMAC256(jwtSecret))
                    call.respond(
                        mapOf("status" to "success", "token" to generatedToken),
                    )
                }
                call.respond(HttpStatusCode.BadRequest, mapOf("status" to "failed"))
            }
            authenticate("auth-jwt") {
                get("me") {
                    call.safeJwt { principal ->
                        call.respond(
                            mapOf(
                                "status" to "success",
                                "id" to principal.payload.getClaim("userid").toString(),
                                "expiry" to principal.expiresAt.toString(),
                            ),
                        )
                    }
                }
            }
            get("users") {
                call.respond(
                    authService.getUsers(),
                )
            }
        }
    }
}
