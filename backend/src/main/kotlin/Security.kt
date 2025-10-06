package net.kazugmx.acadule

import com.auth0.jwt.JWT
import com.auth0.jwt.algorithms.Algorithm
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
import org.jetbrains.exposed.sql.Database
import java.util.*

fun Application.configureSecurity(database: Database) {
    // Please read the jwt property from the config file if you are using EngineMain
    val jwtIssuer = environment.config.property("jwt.issuer").getString()
    val jwtAudience = environment.config.property("jwt.audience").getString()
    val jwtDomain = environment.config.property("jwt.domain").getString()
    val jwtRealm = environment.config.property("jwt.realm").getString()
    val jwtSecret = environment.config.property("jwt.secret").getString()


    authentication {
        jwt {
            realm = jwtRealm
            verifier(
                JWT
                    .require(Algorithm.HMAC256(jwtSecret))
                    .withAudience(jwtAudience)
                    .withIssuer(jwtDomain)
                    .build()
            )
            validate { credential ->
                if (credential.payload.audience.contains(jwtAudience)) JWTPrincipal(credential.payload) else null
            }
        }
    }


    val authService = AuthService(database)
    routing {
        route("/auth") {
            post("createUser") {
                val user = call.receive<UserCreateReq>()
                val result = authService.createUser(user)
                if(result.status)
                    call.respond(result)
                else
                    call.respond(HttpStatusCode.BadRequest,mapOf("status" to "failed"))
            }
            post("login") {
                val login = call.receive<LoginReq>()
                val chall = authService.login(login)
                if (chall != null) {
                    call.respond(
                        JWT.create()
                            .withAudience(jwtAudience)
                            .withIssuer(jwtIssuer)
                            .withClaim("userid", chall.userID)
                            .withExpiresAt(Date(System.currentTimeMillis() + 60000))
                            .sign(Algorithm.HMAC256(jwtSecret))
                    )
                }
                call.respond(mapOf("status" to "failed"))
            }
            get("users") {
                call.respond(
                    authService.getUsers()
                )
            }
        }
    }
}
