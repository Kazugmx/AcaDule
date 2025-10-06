package net.kazugmx.acadule.schemas

import kotlinx.coroutines.Dispatchers
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.Database
import org.jetbrains.exposed.sql.SchemaUtils
import org.jetbrains.exposed.sql.kotlin.datetime.datetime
import org.jetbrains.exposed.sql.transactions.experimental.newSuspendedTransaction
import org.jetbrains.exposed.sql.transactions.transaction
import kotlinx.datetime.LocalDateTime
import org.jetbrains.exposed.sql.insert

@Serializable
enum class TaskProgress {
    NOT_STARTED, IN_PROGRESS, COMPLETE, SUSPENDED
}

@Serializable
data class TaskRes(
    val id: String,
    val ownerID: Int? = null,
    val targetName: String,
    val description: String = "none",
    val progress: TaskProgress = TaskProgress.NOT_STARTED,
    val deadline: LocalDateTime? = null,
    val hasDone: Boolean = false
)

@Serializable
data class CreateTaskReq(
    val ownerID: Int? = null,
    val targetName: String,
    val description: String = "none",
    val progress: TaskProgress = TaskProgress.NOT_STARTED,
    val deadline: LocalDateTime? = null,
)


class TaskService(@Suppress("unused") database: Database) {
    object TasksTable : UUIDTable("tasks") {
        val ownerID = integer("ownerID").nullable()
        val targetName = varchar("tasktarget", 100).default("Untitled")
        val description = text("description").default("none")
        val progress = enumerationByName<TaskProgress>("progress", 64)
            .default(TaskProgress.NOT_STARTED)
        val deadline = datetime("deadline").nullable()
        val hasDone = bool("hasDone").default(false)
    }

    init {
        transaction {
            SchemaUtils.create(TasksTable)
        }
    }

    suspend fun createTask(task: CreateTaskReq) = dbQuery {
        TasksTable.insert {
            it[ownerID] = task.ownerID
            it[targetName] = task.targetName
            it[description] = task.description
            it[progress] = task.progress
            it[deadline] = task.deadline
        } get TasksTable.id
    }

    suspend fun getAllTasksByUser(userID: Int?): List<TaskRes> = dbQuery {
        TasksTable.select(
            TasksTable.id,
            TasksTable.description,
            TasksTable.progress,
            TasksTable.deadline,
            TasksTable.hasDone
        )
            .where { TasksTable.ownerID eq userID }.map {
                TaskRes(
                    id = it[TasksTable.id].value.toString(),
                    targetName = it[TasksTable.targetName],
                    description = it[TasksTable.description],
                    progress = it[TasksTable.progress],
                    deadline = it[TasksTable.deadline],
                    hasDone = it[TasksTable.hasDone]
                )
            }
    }

    private suspend fun <T> dbQuery(block: suspend () -> T): T =
        newSuspendedTransaction(Dispatchers.IO) { block() }
}
