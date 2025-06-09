<script setup>
import { RouterView } from "vue-router";
</script>

<template>
  <div id="app">
    <el-container class="app-container">
      <el-header>
        <div class="header-content">
          <h1>Task Manager</h1>
          <el-menu
            mode="horizontal"
            :default-active="$route.path"
            router
            class="header-menu"
          >
            <el-menu-item index="/tasks">
              <el-icon><List /></el-icon>
              <span>Tasks</span>
            </el-menu-item>
            <el-menu-item index="/tasks/create">
              <el-icon><Plus /></el-icon>
              <span>Create Task</span>
            </el-menu-item>
            <el-menu-item index="/settings">
              <el-icon><Setting /></el-icon>
              <span>Settings</span>
            </el-menu-item>
          </el-menu>
        </div>
      </el-header>
      <el-main>
        <div class="content-container">
          <RouterView />
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<style>
html,
body {
  margin: 0;
  padding: 0;
  width: 100%;
  height: 100%;
}

#app {
  min-height: 100vh;
  width: 100vw;
  margin: 0;
  padding: 0;
}

.app-container {
  min-height: 100vh;
  width: 100%;
}

.el-header {
  background-color: #409eff;
  color: white;
  padding: 0;
  height: 60px;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 100%;
  padding: 0 20px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.header-content h1 {
  margin: 0;
  color: white;
}

.header-menu {
  background-color: transparent;
  border-bottom: none;
}

.header-menu .el-menu-item {
  color: white;
  border-bottom: 2px solid transparent;
}

.header-menu .el-menu-item:hover,
.header-menu .el-menu-item.is-active {
  background-color: rgba(255, 255, 255, 0.1);
  border-bottom-color: white;
  color: white;
}

.el-main {
  padding: 20px;
  background-color: #f5f5f5;
  min-height: calc(100vh - 60px);
}

.content-container {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 20px;
}
</style>
