package net.kazugmx.acadule.schemas

import kotlinx.coroutines.Dispatchers
import kotlinx.datetime.LocalDateTime
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.kotlin.datetime.CurrentDateTime
import org.jetbrains.exposed.sql.kotlin.datetime.datetime
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction
import org.jetbrains.exposed.sql.update
import java.util.*

@Serializable
@Suppress("unused")
enum class TaskProgress {
    NOT_STARTED,
    IN_PROGRESS,
    COMPLETE,
    SUSPENDED,
}

@Serializable
data class TaskRes(
    val id: String,
    val ownerID: Int? = null,
    val title: String,
    val description: String = "none",
    val progress: TaskProgress = TaskProgress.NOT_STARTED,
    val deadline: LocalDateTime? = null,
    val lastUpdated: LocalDateTime,
    val hasDone: Boolean = false,
)

@Serializable
data class CreateTaskReq(
    val title: String,
    val description: String = "none",
    val progress: TaskProgress = TaskProgress.NOT_STARTED,
    val deadline: LocalDateTime? = null,
)

@Serializable
data class UpdateTaskReq(
    val id: String,
    val title: String? = null,
    val description: String? = null,
    val progress: TaskProgress? = TaskProgress.NOT_STARTED,
    val deadline: LocalDateTime? = null,
    val hasDone: Boolean? = null,
    val lastUpdated: LocalDateTime? = null,
)

@Serializable
data class IDTaskReq(
    val id: String,
)

class TaskService(
    @Suppress("unused") database: Database,
) {
    object TasksTable : UUIDTable("tasks") {
        val ownerID = integer("ownerID").nullable()
        val title = varchar("title", 100).default("Untitled")
        val description = text("description").default("none")
        val progress =
            enumerationByName<TaskProgress>("progress", 64)
                .default(TaskProgress.NOT_STARTED)
        val deadline = datetime("deadline").nullable()
        val lastUpdated = datetime("lastUpdated").defaultExpression(CurrentDateTime)
        val hasDone = bool("hasDone").default(false)
    }

    init {
        transaction {
            SchemaUtils.create(TasksTable)
        }
    }

    suspend fun createTask(
        task: CreateTaskReq,
        reqOwnerID: Int,
    ) = dbQuery {
        TasksTable.insert {
            it[ownerID] = reqOwnerID
            it[title] = task.title
            it[description] = task.description
            it[progress] = task.progress
            it[deadline] = task.deadline
        } get TasksTable.id
    }

    suspend fun getAllTasksByUser(userID: Int?): List<TaskRes> =
        dbQuery {
            TasksTable
                .select(
                    TasksTable.id,
                    TasksTable.title,
                    TasksTable.description,
                    TasksTable.progress,
                    TasksTable.deadline,
                    TasksTable.lastUpdated,
                    TasksTable.hasDone,
                ).where { TasksTable.ownerID eq userID }
                .map {
                    TaskRes(
                        id = it[TasksTable.id].value.toString(),
                        title = it[TasksTable.title],
                        description = it[TasksTable.description],
                        progress = it[TasksTable.progress],
                        deadline = it[TasksTable.deadline],
                        lastUpdated = it[TasksTable.lastUpdated],
                        hasDone = it[TasksTable.hasDone],
                    )
                }
        }

    suspend fun getTaskByID(
        taskID: String,
        userID: Int?,
    ): TaskRes? =
        dbQuery {
            TasksTable
                .select(
                    TasksTable.id,
                    TasksTable.title,
                    TasksTable.description,
                    TasksTable.progress,
                    TasksTable.deadline,
                    TasksTable.lastUpdated,
                    TasksTable.hasDone,
                ).limit(1)
                .where {
                    (TasksTable.id eq UUID.fromString(taskID)) and (TasksTable.ownerID eq userID)
                }.map {
                    TaskRes(
                        id = it[TasksTable.id].value.toString(),
                        title = it[TasksTable.title],
                        description = it[TasksTable.description],
                        progress = it[TasksTable.progress],
                        deadline = it[TasksTable.deadline],
                        lastUpdated = it[TasksTable.lastUpdated],
                        hasDone = it[TasksTable.hasDone],
                    )
                }.singleOrNull()
        }

    suspend fun updateTask(
        task: UpdateTaskReq,
        userID: Int,
    ) = dbQuery {
        TasksTable.update(
            where = { (TasksTable.id eq UUID.fromString(task.id)) and (TasksTable.ownerID eq userID) },
        ) {
            if (task.title != null) it[title] = task.title
            if (task.description != null) it[description] = task.description
            if (task.progress != null) it[progress] = task.progress
            if (task.deadline != null) it[deadline] = task.deadline
            if (task.hasDone != null) it[hasDone] = task.hasDone
            it[lastUpdated] = CurrentDateTime
        }
        return@dbQuery getTaskByID(task.id, userID)
    }

    suspend fun deleteTask(
        id: String,
        userID: Int,
    ) = dbQuery {
        TasksTable.deleteWhere { (TasksTable.ownerID eq userID) and (TasksTable.id eq UUID.fromString(id)) }
    }

    private suspend fun <T> dbQuery(block: suspend () -> T): T = newSuspendedTransaction(Dispatchers.IO) { block() }
}
