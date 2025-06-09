<template>
  <div class="task-list">
    <div class="card-header">
      <span>Task List</span>
      <el-button type="primary" @click="$router.push('/tasks/create')">
        <el-icon><Plus /></el-icon>
        Create Task
      </el-button>
    </div>

    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      :closable="false"
      style="margin-bottom: 20px"
    />

    <el-table
      v-loading="loading"
      :data="tasks"
      style="width: 100%"
      @row-click="handleRowClick"
      class="task-table"
    >
      <el-table-column prop="task.title" label="Title" min-width="200">
        <template #default="{ row }">
          <el-link
            type="primary"
            @click="$router.push(`/tasks/${row.task.id}`)"
          >
            {{ row.task.title }}
          </el-link>
        </template>
      </el-table-column>

      <el-table-column
        prop="task.description"
        label="Description"
        min-width="300"
      >
        <template #default="{ row }">
          <el-text class="description-text" truncated>
            {{ row.task.description }}
          </el-text>
        </template>
      </el-table-column>

      <el-table-column prop="task.priority" label="Priority" width="120">
        <template #default="{ row }">
          <el-tag :type="getPriorityColor(row.task.priority)">
            {{ getPriorityText(row.task.priority) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="task.status" label="Status" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusColor(row.task.status)">
            {{ getStatusText(row.task.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="urgency_score" label="Urgency" width="100">
        <template #default="{ row }">
          <el-progress
            :percentage="Math.round(row.urgency_score)"
            :color="getUrgencyColor(row.urgency_score)"
            :stroke-width="8"
          />
        </template>
      </el-table-column>

      <el-table-column
        prop="time_remaining"
        label="Time Remaining"
        width="150"
      >
        <template #default="{ row }">
          <el-text :type="getTimeRemainingType(row.time_remaining)">
            {{ row.time_remaining }}
          </el-text>
        </template>
      </el-table-column>

      <el-table-column prop="task.deadline" label="Deadline" width="180">
        <template #default="{ row }">
          {{ formatDate(row.task.deadline) }}
        </template>
      </el-table-column>

      <el-table-column label="Actions" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button
              size="small"
              type="primary"
              @click.stop="updateStatus(row.task.id, 1)"
              :disabled="row.task.status === 1"
            >
              Start
            </el-button>
            <el-button
              size="small"
              type="success"
              @click.stop="updateStatus(row.task.id, 2)"
              :disabled="row.task.status === 2"
            >
              Complete
            </el-button>
            <el-button
              size="small"
              type="danger"
              @click.stop="deleteTask(row.task.id)"
            >
              Delete
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="!loading && tasks.length === 0" class="empty-state">
      <el-empty description="No tasks found">
        <el-button type="primary" @click="$router.push('/tasks/create')">
          Create Your First Task
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useTaskStore } from "../stores/task";
import { ElMessage, ElMessageBox } from "element-plus";

const taskStore = useTaskStore();
const {
  tasks,
  loading,
  error,
  fetchTasks,
  updateTaskStatus,
  deleteTask: deleteTaskFromStore,
  getPriorityColor,
  getPriorityText,
  getStatusText,
  getStatusColor,
} = taskStore;

onMounted(() => {
  fetchTasks();
});

const handleRowClick = (row) => {
  // Navigate to task detail
};

const updateStatus = async (taskId, status) => {
  try {
    await updateTaskStatus(taskId, status);
    ElMessage.success("Task status updated successfully");
  } catch (error) {
    ElMessage.error("Failed to update task status");
  }
};

const deleteTask = async (taskId) => {
  try {
    await ElMessageBox.confirm(
      "This will permanently delete the task. Continue?",
      "Warning",
      {
        confirmButtonText: "OK",
        cancelButtonText: "Cancel",
        type: "warning",
      },
    );
    await deleteTaskFromStore(taskId);
    ElMessage.success("Task deleted successfully");
  } catch (error) {
    if (error !== "cancel") {
      ElMessage.error("Failed to delete task");
    }
  }
};

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleString();
};

const getUrgencyColor = (score) => {
  if (score >= 80) return "#f56c6c";
  if (score >= 60) return "#e6a23c";
  if (score >= 40) return "#409eff";
  return "#67c23a";
};

const getTimeRemainingType = (timeRemaining) => {
  if (timeRemaining === "Overdue") return "danger";
  if (timeRemaining.includes("hours") && !timeRemaining.includes("days"))
    return "warning";
  return "info";
};
</script>

<style scoped>
.task-list {
  width: 100%;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.card-header span {
  font-size: 18px;
  font-weight: bold;
}

.description-text {
  display: block;
  max-width: 300px;
}

.empty-state {
  margin: 40px 0;
  text-align: center;
}

.task-table {
  margin-bottom: 20px;
}

/* 修复Element Plus组件文字显示问题 */
.el-table {
  font-size: 14px;
  color: #333;
}

.el-table th {
  font-weight: 600;
  background-color: #f5f7fa;
}

.el-card {
  width: 100%;
  margin-bottom: 20px;
}

.el-card__header {
  padding: 15px 20px;
  font-weight: 600;
}

.el-button {
  font-weight: 500;
}

.el-tag {
  font-weight: 500;
}
</style>
