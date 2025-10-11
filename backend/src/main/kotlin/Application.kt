package net.kazugmx.acadule

import com.zaxxer.hikari.HikariConfig
import com.zaxxer.hikari.HikariDataSource
import io.ktor.serialization.kotlinx.json.json
import io.ktor.server.application.*
import io.ktor.server.auth.Authentication
import io.ktor.server.plugins.calllogging.CallLogging
import io.ktor.server.plugins.contentnegotiation.ContentNegotiation
import io.ktor.server.request.path
import kotlinx.serialization.json.Json
import net.kazugmx.acadule.components.configureAuth
import net.kazugmx.acadule.components.configureTaskService
import net.kazugmx.acadule.schemas.AuthService
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.transactions.TransactionManager
import kotlin.system.exitProcess
import kotlin.time.ExperimentalTime

fun main(args: Array<String>) {
    io.ktor.server.netty.EngineMain.main(args)
}

@Suppress("unused")
private object DBDriver {
    const val MYSQL = "org.mariadb.jdbc.Driver"
    const val SQLITE = "org.sqlite.JDBC"
    const val POSTGRES = "org.postgresql.Driver"
}

@OptIn(ExperimentalTime::class)
fun Application.module() {
    val config = HikariConfig().apply {
        jdbcUrl = environment.config.property("db.url").getString()
        driverClassName = "org.postgresql.Driver"
        maximumPoolSize = 10
        isAutoCommit = true
        transactionIsolation = "TRANSACTION_SERIALIZABLE"
        username = environment.config.property("db.user").getString()
        password = environment.config.property("db.password").getString()
        validate()
    }
    val dataSource = HikariDataSource(config)

    val database = Database.connect(datasource = dataSource)


    val authService = AuthService(database)

    install(ContentNegotiation) {
        json(Json {
            prettyPrint = true
            isLenient = true

        })
    }

    install(CallLogging) {
        filter { call ->
            call.request.path().startsWith("/auth")
        }
    }
    install(Authentication)

    configureHTTP()
    configureAuth(authService)
    configureTaskService(database, authService)
    configureRouting()
}
