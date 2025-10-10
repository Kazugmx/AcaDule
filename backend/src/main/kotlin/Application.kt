package net.kazugmx.acadule

import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.*
import io.ktor.server.auth.Authentication
import io.ktor.server.auth.jwt.JWTPrincipal
import io.ktor.server.auth.jwt.jwt
import io.ktor.server.plugins.calllogging.CallLogging
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.request.path
import net.kazugmx.acadule.schemas.AuthService
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.transactions.TransactionManager

fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

fun Application.module() {
    TransactionManager.manager.defaultIsolationLevel = java.sql.Connection.TRANSACTION_SERIALIZABLE
    val database = Database.connect(
        url = "jdbc:sqlite:data.db",
        driver = "org.sqlite.JDBC"
    )
    val authService = AuthService(database)

    install(ContentNegotiation) {
        json()
    }

    install(CallLogging) {
        filter { call ->
            call.request.path().startsWith("/auth")
        }

    }
    install(Authentication)

    configureHTTP()
    configureAuth(authService)
    configureDatabases(database)
    configureRouting()
}
