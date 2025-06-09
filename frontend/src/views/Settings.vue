<script setup>
import { ref, onMounted } from "vue";
import { useSettingsStore } from "../stores/settings";
import { ElMessage } from "element-plus";

const settingsStore = useSettingsStore();
const apiKey = ref("");
const model = ref("");
const modelOptions = [
  { value: "gpt-4", label: "GPT-4" },
  { value: "gpt-4o", label: "GPT-4o" },
  { value: "gpt-4-turbo", label: "GPT-4 Turbo" },
  { value: "gpt-4o-mini", label: "GPT-4o Mini" },
  { value: "gpt-3.5-turbo", label: "GPT-3.5 Turbo" },
];
const loading = ref(false);

onMounted(async () => {
  loading.value = true;
  try {
    await settingsStore.fetchSettings();
    apiKey.value = settingsStore.apiKey || "";
    model.value = settingsStore.model || "gpt-4o-mini";
  } catch (error) {
    ElMessage.error("Failed to load settings");
  } finally {
    loading.value = false;
  }
});

const saveSettings = async () => {
  loading.value = true;
  try {
    await settingsStore.updateSettings({
      api_key: apiKey.value,
      model: model.value,
    });
    ElMessage.success("Settings saved successfully");
  } catch (error) {
    ElMessage.error("Failed to save settings");
  } finally {
    loading.value = false;
  }
};
</script>

<template>
  <div class="settings-container">
    <div class="card-header">
      <span>API Settings</span>
    </div>

    <el-form
      label-position="top"
      :model="{ apiKey, model }"
      v-loading="loading"
    >
      <el-form-item label="OpenAI API Key">
        <el-input
          v-model="apiKey"
          placeholder="sk-..."
          show-password
          clearable
        />
        <div class="form-tip">
          API key must have access to GPT-4 models. You can get your API key
          from
          <a href="https://platform.openai.com/api-keys" target="_blank"
            >OpenAI dashboard</a
          >.
        </div>
      </el-form-item>

      <el-form-item label="Model">
        <el-select v-model="model" class="model-select">
          <el-option
            v-for="option in modelOptions"
            :key="option.value"
            :label="option.label"
            :value="option.value"
          />
        </el-select>
        <div class="form-tip">
          Select GPT-4 model to use for task generation.
        </div>
      </el-form-item>

      <div class="actions">
        <el-button type="primary" @click="saveSettings" :loading="loading">
          Save Settings
        </el-button>
      </div>
    </el-form>
  </div>
</template>

<style scoped>
.settings-container {
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

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.model-select {
  width: 100%;
}

.actions {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
}
</style>
