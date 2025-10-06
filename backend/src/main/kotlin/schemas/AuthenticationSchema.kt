package net.kazugmx.acadule.schemas

import at.favre.lib.crypto.bcrypt.BCrypt
import kotlinx.coroutines.Dispatchers
import kotlinx.datetime.LocalDateTime
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.kotlin.datetime.*
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction


@Serializable
data class UserRes(
    val username: String,
    val mail: String,
    val createdAt: LocalDateTime? = null,
    val lastLoginAt: LocalDateTime? = null,
)

@Serializable
data class UserCreateReq(
    val username: String,
    val mail: String,
    val password: String
)
@Serializable
data class UserCreateRes(
    val status: Boolean,
    val id: Int
)

@Serializable
data class LoginReq(
    val username: String,
    val password: String
)

data class LoginRes(
    val userID: Int
)


class AuthService(database: Database) {
    object UserTable : IntIdTable("users") {
        val mail = varchar("mail", length = 255).uniqueIndex()
        val username = varchar("username", 50).uniqueIndex()
        val password = varchar("password", 100)
        val createdAt = datetime("created_at").defaultExpression(CurrentDateTime)
        val lastLoginAt = datetime("last_login_at").nullable()
    }


    init {
        transaction {
            SchemaUtils.create(UserTable)
        }
    }

    suspend fun createUser(user: UserCreateReq):UserCreateRes = dbQuery {
        val passHash = BCrypt.withDefaults().hashToString(12, user.password.toCharArray())
        try {
            val id: Int = UserTable.insert {
                it[mail] = user.mail
                it[username] = user.username
                it[password] = passHash
            }[UserTable.id].value
            return@dbQuery UserCreateRes(true,id)
        } catch(e: Exception) {
            exposedLogger.error(e.message)
            return@dbQuery UserCreateRes(false,-1)
        }
    }

    suspend fun login(loginReq: LoginReq): LoginRes? = dbQuery {
        val userData = UserTable
            .select(UserTable.id, UserTable.password)
            .limit(1)
            .where { UserTable.username eq loginReq.username }
            .singleOrNull() ?: return@dbQuery null
        val check = BCrypt.verifyer().verify(
            loginReq.password.toCharArray(),
            userData[UserTable.password].toCharArray()
        ).verified
        if (check) {
            UserTable.update({ UserTable.id eq userData[UserTable.id] }) {
                it[lastLoginAt] = CurrentDateTime
            }
            return@dbQuery LoginRes(userData[UserTable.id].value)
        }
        return@dbQuery null
    }

    suspend fun getUsers() = dbQuery {
        UserTable.selectAll().map {
            UserRes(
                username = it[UserTable.username],
                mail = it[UserTable.mail],
                createdAt = it[UserTable.createdAt],
                lastLoginAt = it[UserTable.lastLoginAt]
            )
        }
    }

    suspend fun deleteUser(id: Int) {
        dbQuery {
            UserTable.deleteWhere { UserTable.id eq id }
        }
    }

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction(Dispatchers.IO) { block() }
}