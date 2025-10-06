package net.kazugmx.acadule

import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.*
import io.ktor.server.plugins.calllogging.CallLogging
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.request.path
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

    install(ContentNegotiation) {
        json()
    }

    install(CallLogging){
        filter { call ->
            call.request.path().startsWith("/auth")
        }

    }

    configureHTTP()
    configureSecurity(database)
    configureDatabases(database)
    configureRouting()
}
