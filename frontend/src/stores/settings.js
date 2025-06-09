import { defineStore } from "pinia";
import { ref } from "vue";
import axios from "axios";

const API_BASE_URL = "http://localhost:8080/api/v1";

export const useSettingsStore = defineStore("settings", () => {
  const apiKey = ref("");
  const model = ref("gpt-4o-mini");
  const useLLM = ref(true);
  const loading = ref(false);
  const error = ref(null);

  const api = axios.create({
    baseURL: API_BASE_URL,
    headers: {
      "Content-Type": "application/json",
    },
  });

  const fetchSettings = async () => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.get("/config/openai");
      apiKey.value = response.data.api_key || "";
      model.value = response.data.model || "gpt-4o-mini";
      useLLM.value = response.data.use_llm;
      return response.data;
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to fetch settings";
      console.error("Error fetching settings:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateSettings = async (settings) => {
    loading.value = true;
    error.value = null;
    try {
      const response = await api.put("/config/openai", settings);
      apiKey.value = response.data.api_key || "";
      model.value = response.data.model || "gpt-4o-mini";
      return response.data;
    } catch (err) {
      error.value = err.response?.data?.error || "Failed to update settings";
      console.error("Error updating settings:", err);
      throw err;
    } finally {
      loading.value = false;
    }
  };

  return {
    apiKey,
    model,
    useLLM,
    loading,
    error,
    fetchSettings,
    updateSettings,
  };
});
