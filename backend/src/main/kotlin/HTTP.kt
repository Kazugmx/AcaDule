package net.kazugmx.acadule

import io.ktor.server.application.Application
import io.ktor.server.application.install
import io.ktor.server.plugins.forwardedheaders.ForwardedHeaders
import io.ktor.server.plugins.forwardedheaders.XForwardedHeaders
import io.ktor.server.response.respond
import io.ktor.server.routing.get
import io.ktor.server.routing.routing

fun Application.configureHTTP() {
    install(ForwardedHeaders) // WARNING: for security, do not include this if not behind a reverse proxy
    install(XForwardedHeaders) // WARNING: for security, do not include this if not behind a reverse proxy
    routing {
        get("/health") {
            call.respond(message = "ok")
        }
    }
}
