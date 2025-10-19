package net.kazugmx.acadule.components

import com.auth0.jwt.exceptions.JWTDecodeException
import com.auth0.jwt.exceptions.TokenExpiredException
import io.ktor.http.HttpStatusCode
import io.ktor.server.application.ApplicationCall
import io.ktor.server.application.log
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.principal
import io.ktor.server.response.respond

suspend inline fun ApplicationCall.safeJwt(block: (JWTPrincipal) -> Unit) {
    try {
        val principal =
            principal<JWTPrincipal>() ?: return respond(
                HttpStatusCode.BadRequest,
                mapOf("status" to "failed", "reason" to "No JWT Principal"),
            )
        block(principal)
    } catch (e: TokenExpiredException) {
        respond(
            HttpStatusCode.BadRequest,
            mapOf(
                "status" to "failed",
                "reason" to "Token Expired",
                "expiredOn" to e.expiredOn.toString(),
            ),
        )
        application.log.info("Token Expired: ${e.expiredOn}")
    } catch (e: JWTDecodeException) {
        respond(
            HttpStatusCode.BadRequest,
            mapOf(
                "status" to "failed",
                "reason" to "Invalid Token",
                "detail" to e.message,
            ),
        )
    } catch (e: Exception) {
        application.log.error("JWT processing failed", e)
        respond(
            HttpStatusCode.InternalServerError,
            mapOf(
                "status" to "failed",
                "reason" to "Internal Server Error",
            ),
        )
    }
}
