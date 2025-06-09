<template>
  <div class="task-create">
    <div class="card-header">
      <span>Create New Task</span>
      <el-button @click="$router.push('/tasks')">
        <el-icon><ArrowLeft /></el-icon>
        Back to Tasks
      </el-button>
    </div>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
      @submit.prevent="submitForm"
    >
      <el-form-item label="Title" prop="title">
        <el-input
          v-model="form.title"
          placeholder="Enter task title"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="Description" prop="description">
        <el-input
          v-model="form.description"
          type="textarea"
          :rows="6"
          placeholder="Describe your task in detail..."
          maxlength="1000"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="Deadline" prop="deadline">
        <el-date-picker
          v-model="form.deadline"
          type="datetime"
          placeholder="Select deadline"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DDTHH:mm:ss.SSSZ"
          :disabled-date="disabledDate"
          style="width: 100%"
        />
      </el-form-item>

      <el-form-item>
        <el-button
          type="primary"
          @click="submitForm"
          :loading="loading"
          size="large"
        >
          <el-icon><Check /></el-icon>
          Create Task
        </el-button>
        <el-button @click="resetForm" size="large">
          <el-icon><Refresh /></el-icon>
          Reset
        </el-button>
      </el-form-item>
    </el-form>

    <el-alert
      v-if="error"
      :title="error"
      type="error"
      show-icon
      :closable="false"
      style="margin-top: 20px"
    />

    <el-alert
      v-if="isOpenAIError"
      title="API Key Configuration Required"
      type="warning"
      show-icon
      :closable="false"
      style="margin-top: 10px"
    >
      <div>
        Please configure your OpenAI API key with GPT-4 access in the
        <router-link to="/settings">Settings</router-link> page.
      </div>
    </el-alert>

    <!-- Preview Section -->
    <div v-if="form.title || form.description" class="task-preview">
      <h3>Task Preview</h3>
      <div class="preview-content">
        <h3 v-if="form.title">{{ form.title }}</h3>
        <p v-if="form.description" class="description">
          {{ form.description }}
        </p>
        <div v-if="form.deadline" class="deadline-info">
          <el-tag type="info">
            <el-icon><Clock /></el-icon>
            Deadline: {{ formatDeadline(form.deadline) }}
          </el-tag>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useTaskStore } from "../stores/task";
import { ElMessage } from "element-plus";

const router = useRouter();
const taskStore = useTaskStore();
const { createTask, loading, error, isOpenAIError } = taskStore;

const formRef = ref();

const form = reactive({
  title: "",
  description: "",
  deadline: null,
});

const rules = {
  title: [
    { required: true, message: "Please enter task title", trigger: "blur" },
    {
      min: 3,
      max: 100,
      message: "Title should be 3-100 characters",
      trigger: "blur",
    },
  ],
  description: [
    {
      required: true,
      message: "Please enter task description",
      trigger: "blur",
    },
    {
      min: 10,
      max: 1000,
      message: "Description should be 10-1000 characters",
      trigger: "blur",
    },
  ],
  deadline: [
    { required: true, message: "Please select deadline", trigger: "change" },
  ],
};

const submitForm = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();

    const taskData = {
      title: form.title,
      description: form.description,
      deadline: form.deadline,
    };

    await createTask(taskData);
    ElMessage.success("Task created successfully!");
    router.push("/tasks");
  } catch (error) {
    if (error !== false) {
      // validation error returns false
      const errorMessage = taskStore.error || "Failed to create task";
      ElMessage.error(errorMessage);
    }
  }
};

const resetForm = () => {
  if (formRef.value) {
    formRef.value.resetFields();
  }
  Object.assign(form, {
    title: "",
    description: "",
    deadline: null,
  });
};

const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7; // Disable past dates
};

const formatDeadline = (deadline) => {
  if (!deadline) return "";
  return new Date(deadline).toLocaleString();
};
</script>

<style scoped>
.task-create {
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

.task-preview {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #ebeef5;
}

.task-preview h3 {
  margin-top: 0;
  margin-bottom: 20px;
  font-size: 16px;
  font-weight: bold;
}

.preview-content {
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 4px;
}

.preview-content h3 {
  margin-top: 0;
}

.description {
  white-space: pre-wrap;
  margin-bottom: 15px;
}

.deadline-info {
  margin-top: 15px;
}
</style>
