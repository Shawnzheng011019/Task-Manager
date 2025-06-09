import { defineStore } from "pinia";
import { ref, computed } from "vue";
import axios from "axios";

const API_BASE_URL = "http://localhost:8080/api/v1";

export const useTaskStore = defineStore("task", () => {
  const tasks = ref([]);
  const currentTask = ref(null);
  const loading = ref(false);
  const error = ref(null);

  const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
      "Content-Type": "application/json",
    },
  });

  const isOpenAIError = computed(() => {
    if (!error.value) return false;
    return (
      error.value.includes("OpenAI API key is invalid") ||
      error.value.includes("OpenAI API rate limit") ||
      error.value.includes("OpenAI API request is invalid") ||
      error.value.includes("model availability")
    );
  });

  const fetchTasks = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.get("/tasks");
      tasks.value = response.data.tasks || [];
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to fetch tasks";
      console.error("Error fetching tasks:", err);
    } finally {
      loading.value = false;
    }
  };

  const fetchTaskById = async (id) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.get(`/tasks/${id}`);
      currentTask.value = response.data;
      return response.data;
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to fetch task";
      console.error("Error fetching task:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const createTask = async (taskData) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.post("/tasks", taskData);
      await fetchTasks(); // Refresh the task list
      return response.data;
    } catch (err) {
      if (
        err.response?.data?.error &&
        err.response.data.error.includes("OpenAI API key is invalid")
      ) {
        error.value =
          "OpenAI API key is invalid or not configured. Please check your server configuration.";
      } else if (
        err.response?.data?.error &&
        err.response.data.error.includes("OpenAI API rate limit")
      ) {
        error.value = "OpenAI API rate limit reached. Please try again later.";
      } else if (
        err.response?.data?.error &&
        err.response.data.error.includes("OpenAI API request is invalid")
      ) {
        error.value =
          "Unable to connect to AI service. Please check your configuration or try again later.";
      } else {
        error.value = err.response?.data?.error || "Failed to create task";
      }
      console.error("Error creating task:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateTaskStatus = async (id, status) => {
    loading.value = true;
    error.value = null;
    try {
      await api.put(`/tasks/${id}/status`, { status });
      await fetchTasks(); // Refresh the task list
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to update task status";
      console.error("Error updating task status:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteTask = async (id) => {
    loading.value = true;
    error.value = null;
    try {
      await api.delete(`/tasks/${id}`);
      await fetchTasks(); // Refresh the task list
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to delete task";
      console.error("Error deleting task:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const getPriorityColor = (priority) => {
    switch (priority) {
      case 0:
        return "info"; // Low
      case 1:
        return "warning"; // Medium
      case 2:
        return "danger"; // High
      case 3:
        return "danger"; // Urgent
      default:
        return "info";
    }
  };

  const getPriorityText = (priority) => {
    switch (priority) {
      case 0:
        return "Low";
      case 1:
        return "Medium";
      case 2:
        return "High";
      case 3:
        return "Urgent";
      default:
        return "Unknown";
    }
  };

  const getStatusText = (status) => {
    switch (status) {
      case 0:
        return "Pending";
      case 1:
        return "In Progress";
      case 2:
        return "Completed";
      case 3:
        return "Cancelled";
      default:
        return "Unknown";
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 0:
        return "warning"; // Pending
      case 1:
        return "primary"; // In Progress
      case 2:
        return "success"; // Completed
      case 3:
        return "info"; // Cancelled
      default:
        return "info";
    }
  };

  return {
    tasks,
    currentTask,
    loading,
    error,
    isOpenAIError,
    fetchTasks,
    fetchTaskById,
    createTask,
    updateTaskStatus,
    deleteTask,
    getPriorityColor,
    getPriorityText,
    getStatusText,
    getStatusColor,
  };
});
