<template>
  <el-container style="min-height: 100vh;">
    <!-- 侧边导航栏 -->
    <el-aside width="220px" style="background: #304156;">
      <div style="height:60px;display:flex;align-items:center;justify-content:center;color:#fff;font-size:18px;font-weight:bold;border-bottom:1px solid rgba(255,255,255,.1);">
        图书管理系统
      </div>
      <el-menu
        :default-active="route.path"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
        style="border-right: none;"
      >
        <el-menu-item index="/books">
          <el-icon><Reading /></el-icon><span>图书浏览</span>
        </el-menu-item>
        <el-menu-item index="/books/manage" v-if="auth.user?.role === 'admin'">
          <el-icon><Edit /></el-icon><span>图书管理</span>
        </el-menu-item>
        <el-menu-item index="/borrows">
          <el-icon><Tickets /></el-icon><span>我的借阅</span>
        </el-menu-item>
        <el-menu-item index="/profile">
          <el-icon><User /></el-icon><span>个人中心</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container>
      <!-- 顶部栏 -->
      <el-header style="background:#fff;border-bottom:1px solid #e6e6e6;display:flex;align-items:center;justify-content:flex-end;padding:0 20px;">
        <el-dropdown @command="handleCommand">
          <span style="cursor:pointer;display:flex;align-items:center;gap:8px;">
            {{ auth.user?.name || '用户' }}
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人中心</el-dropdown-item>
              <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-header>

      <!-- 页面主体 -->
      <el-main style="background:#f0f2f5;padding:20px;">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

function handleCommand(command) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    ElMessageBox.confirm('确定要退出登录吗？', '提示')
      .then(() => {
        auth.logout()
        router.push('/login')
        ElMessage.success('已退出登录')
      })
      .catch(() => {})
  }
}
</script>
