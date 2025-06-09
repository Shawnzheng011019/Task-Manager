<template>
  <div class="task-detail" v-loading="loading">
    <div class="card-header">
      <span>Task Details</span>
      <el-button @click="$router.push('/tasks')">
        <el-icon><ArrowLeft /></el-icon>
        Back to Tasks
      </el-button>
    </div>

    <div v-if="currentTask" class="task-content">
      <!-- Task Basic Info -->
      <div class="task-header">
        <h2>{{ currentTask.task.title }}</h2>
        <div class="task-meta">
          <el-tag
            :type="getPriorityColor(currentTask.task.priority)"
            size="large"
          >
            {{ getPriorityText(currentTask.task.priority) }} Priority
          </el-tag>
          <el-tag
            :type="getStatusColor(currentTask.task.status)"
            size="large"
          >
            {{ getStatusText(currentTask.task.status) }}
          </el-tag>
        </div>
      </div>

      <!-- Urgency and Time Info -->
      <el-row :gutter="20" style="margin: 20px 0">
        <el-col :span="8">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-label">Urgency Score</div>
              <el-progress
                :percentage="Math.round(currentTask.urgency_score)"
                :color="getUrgencyColor(currentTask.urgency_score)"
                :stroke-width="12"
              />
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-label">Time Remaining</div>
              <div
                class="stat-value"
                :class="getTimeRemainingClass(currentTask.time_remaining)"
              >
                {{ currentTask.time_remaining }}
              </div>
            </div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card shadow="hover">
            <div class="stat-item">
              <div class="stat-label">Deadline</div>
              <div class="stat-value">
                {{ formatDate(currentTask.task.deadline) }}
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- Task Description -->
      <div class="section-card">
        <h3>Description</h3>
        <div class="description-content">
          {{ currentTask.task.description }}
        </div>
      </div>

      <!-- Technical Plan -->
      <div
        v-if="currentTask.task.technical_plan"
        class="section-card"
      >
        <h3>Technical Plan</h3>
        <div
          class="markdown-content"
          v-html="formatMarkdown(currentTask.task.technical_plan)"
        ></div>
      </div>

      <!-- Workflow -->
      <div
        v-if="currentTask.task.workflow"
        class="section-card"
      >
        <h3>Workflow</h3>
        <div
          class="markdown-content"
          v-html="formatMarkdown(currentTask.task.workflow)"
        ></div>
      </div>

      <!-- Sub Tasks -->
      <div
        v-if="
          currentTask.task.sub_tasks && currentTask.task.sub_tasks.length > 0
        "
        class="section-card"
      >
        <h3>Sub Tasks</h3>
        <el-table :data="currentTask.task.sub_tasks" style="width: 100%">
          <el-table-column prop="order" label="#" width="60" />
          <el-table-column prop="title" label="Title" min-width="200" />
          <el-table-column
            prop="description"
            label="Description"
            min-width="300"
          />
          <el-table-column
            prop="estimated_hours"
            label="Est. Hours"
            width="100"
          />
          <el-table-column prop="priority" label="Priority" width="100">
            <template #default="{ row }">
              <el-tag :type="getPriorityColor(row.priority)" size="small">
                {{ getPriorityText(row.priority) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="status" label="Status" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusColor(row.status)" size="small">
                {{ getStatusText(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- Action Buttons -->
      <div class="action-buttons">
        <el-button
          type="primary"
          size="large"
          @click="updateStatus(1)"
          :disabled="currentTask.task.status === 1"
        >
          <el-icon><VideoPlay /></el-icon>
          Start Task
        </el-button>
        <el-button
          type="success"
          size="large"
          @click="updateStatus(2)"
          :disabled="currentTask.task.status === 2"
        >
          <el-icon><Check /></el-icon>
          Mark Complete
        </el-button>
        <el-button
          type="warning"
          size="large"
          @click="updateStatus(3)"
          :disabled="currentTask.task.status === 3"
        >
          <el-icon><Close /></el-icon>
          Cancel Task
        </el-button>
      </div>
    </div>

    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      :closable="false"
      style="margin-top: 20px"
    />
  </div>
</template>

<script setup>
import { onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useTaskStore } from "../stores/task";
import { ElMessage } from "element-plus";

const route = useRoute();
const router = useRouter();
const taskStore = useTaskStore();
const {
  currentTask,
  loading,
  error,
  fetchTaskById,
  updateTaskStatus,
  getPriorityColor,
  getPriorityText,
  getStatusText,
  getStatusColor,
} = taskStore;

const props = defineProps({
  id: {
    type: String,
    required: true,
  },
});

onMounted(async () => {
  try {
    await fetchTaskById(props.id);
  } catch (error) {
    ElMessage.error("Task not found");
    router.push("/tasks");
  }
});

const updateStatus = async (status) => {
  try {
    await updateTaskStatus(currentTask.value.task.id, status);
    await fetchTaskById(props.id); // Refresh task data
    ElMessage.success("Task status updated successfully");
  } catch (error) {
    ElMessage.error("Failed to update task status");
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

const getTimeRemainingClass = (timeRemaining) => {
  if (timeRemaining === "Overdue") return "overdue";
  if (timeRemaining.includes("hours") && !timeRemaining.includes("days"))
    return "urgent";
  return "normal";
};

const formatMarkdown = (text) => {
  // Simple markdown to HTML conversion
  return text
    .replace(/^# (.*$)/gim, "<h1>$1</h1>")
    .replace(/^## (.*$)/gim, "<h2>$1</h2>")
    .replace(/^### (.*$)/gim, "<h3>$1</h3>")
    .replace(/^\* (.*$)/gim, "<li>$1</li>")
    .replace(/^\d+\. (.*$)/gim, "<li>$1</li>")
    .replace(/\*\*(.*)\*\*/gim, "<strong>$1</strong>")
    .replace(/\*(.*)\*/gim, "<em>$1</em>")
    .replace(/\n/gim, "<br>");
};
</script>

<style scoped>
.task-detail {
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

.task-header {
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-bottom: 20px;
}

.task-header h2 {
  margin: 0;
  font-size: 24px;
}

.task-meta {
  display: flex;
  gap: 10px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  font-weight: bold;
  margin-bottom: 10px;
  color: #606266;
}

.stat-value {
  font-size: 18px;
  font-weight: bold;
}

.section-card {
  background-color: white;
  border-radius: 8px;
  padding: 20px;
  margin: 20px 0;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.section-card h3 {
  margin-top: 0;
  margin-bottom: 20px;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #ebeef5;
  padding-bottom: 10px;
}

.description-content {
  white-space: pre-wrap;
  line-height: 1.6;
}

.markdown-content {
  line-height: 1.6;
}

.action-buttons {
  display: flex;
  gap: 10px;
  margin-top: 30px;
  justify-content: center;
}

.past-deadline {
  color: #f56c6c;
}

.approaching-deadline {
  color: #e6a23c;
}

.comfortable-deadline {
  color: #67c23a;
}
</style>
