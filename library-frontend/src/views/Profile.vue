<template>
  <el-card shadow="never">
    <template #header>
      <span>个人信息</span>
    </template>
    <el-descriptions :column="1" border v-if="auth.user">
      <el-descriptions-item label="用户ID">{{ auth.user.id }}</el-descriptions-item>
      <el-descriptions-item label="姓名">{{ auth.user.name }}</el-descriptions-item>
      <el-descriptions-item label="邮箱">{{ auth.user.email }}</el-descriptions-item>
      <el-descriptions-item label="手机号">{{ auth.user.phone || '未绑定' }}</el-descriptions-item>
      <el-descriptions-item label="角色">
        <el-tag :type="auth.user.role === 'admin' ? 'danger' : 'info'" size="small">
          {{ auth.user.role === 'admin' ? '管理员' : '普通用户' }}
        </el-tag>
      </el-descriptions-item>
    </el-descriptions>
    <div v-else style="text-align:center;padding:40px;color:#999;">
      加载中...
    </div>
  </el-card>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()

onMounted(async () => {
  try {
    await auth.fetchProfile()
  } catch (e) {
    // 已在拦截器中处理
  }
})
</script>
