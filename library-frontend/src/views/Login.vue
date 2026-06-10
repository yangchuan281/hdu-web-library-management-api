<template>
  <div class="login-container">
    <div class="login-card">
      <h2 style="text-align:center;margin-bottom:24px;color:#303133;">
        {{ mode === 'login' ? '登录' : mode === 'register' ? '邮箱注册' : '重置密码' }}
      </h2>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="0"
        size="large"
        @keyup.enter="handleSubmit"
      >
        <!-- 邮箱 -->
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="邮箱地址">
            <template #prefix><el-icon><Message /></el-icon></template>
          </el-input>
        </el-form-item>

        <!-- 验证码（注册/重置模式） -->
        <el-form-item v-if="mode !== 'login'" prop="code">
          <div style="display:flex;width:100%;">
            <el-input v-model="form.code" placeholder="4位验证码" maxlength="4">
              <template #prefix><el-icon><Key /></el-icon></template>
            </el-input>
            <el-button :disabled="cooldown > 0" style="width:110px;margin-left:10px;flex-shrink:0;" @click="sendCode">
              {{ cooldown > 0 ? `${cooldown}s` : '获取验证码' }}
            </el-button>
          </div>
        </el-form-item>

        <!-- 密码 -->
        <el-form-item v-if="mode !== 'reset'" prop="password">
          <el-input v-model="form.password" type="password" show-password placeholder="密码">
            <template #prefix><el-icon><Lock /></el-icon></template>
          </el-input>
        </el-form-item>

        <!-- 新密码（重置模式） -->
        <el-form-item v-if="mode === 'reset'" prop="newPassword">
          <el-input v-model="form.newPassword" type="password" show-password placeholder="新密码（至少6位）">
            <template #prefix><el-icon><Lock /></el-icon></template>
          </el-input>
        </el-form-item>

        <!-- 姓名（注册模式） -->
        <el-form-item v-if="mode === 'register'" prop="name">
          <el-input v-model="form.name" placeholder="用户姓名">
            <template #prefix><el-icon><User /></el-icon></template>
          </el-input>
        </el-form-item>

        <!-- 提交按钮 -->
        <el-form-item>
          <el-button type="primary" style="width:100%;" size="large" :loading="loading" @click="handleSubmit">
            {{ mode === 'login' ? '登 录' : mode === 'register' ? '注 册' : '重 置 密 码' }}
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 模式切换 -->
      <div style="text-align:center;font-size:14px;">
        <template v-if="mode === 'login'">
          没有账号？<a href="#" @click.prevent="switchMode('register')">邮箱注册</a>
          &nbsp;|&nbsp;
          <a href="#" @click.prevent="switchMode('reset')">忘记密码？</a>
        </template>
        <template v-else>
          <a href="#" @click.prevent="switchMode('login')">返回登录</a>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { sendCodeApi, registerApi, resetPasswordApi } from '../api/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const auth = useAuthStore()
const formRef = ref(null)
const loading = ref(false)
const mode = ref('login') // login | register | reset
const cooldown = ref(0)

const form = reactive({
  email: '',
  code: '',
  password: '',
  newPassword: '',
  name: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入验证码', trigger: 'blur' },
    { len: 4, message: '验证码为4位数字', trigger: 'blur' },
    { pattern: /^\d{4}$/, message: '验证码必须为数字', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ]
}

function switchMode(m) {
  mode.value = m
  form.code = ''
  form.newPassword = ''
  form.name = ''
}

async function sendCode() {
  if (!form.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  const type = mode.value === 'register' ? 'register' : 'reset'
  try {
    await sendCodeApi(form.email, type)
    ElMessage.success('验证码已发送（请查看服务器日志或数据库）')
    cooldown.value = 60
    const timer = setInterval(() => {
      cooldown.value--
      if (cooldown.value <= 0) clearInterval(timer)
    }, 1000)
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

async function handleSubmit() {
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    if (mode.value === 'login') {
      await auth.login({ email: form.email, password: form.password })
      await auth.fetchProfile()
      ElMessage.success('登录成功')
      router.push('/books')
    } else if (mode.value === 'register') {
      await registerApi({
        email: form.email,
        code: form.code,
        password: form.password,
        name: form.name
      })
      ElMessage.success('注册成功，请登录')
      switchMode('login')
    } else if (mode.value === 'reset') {
      await resetPasswordApi({
        email: form.email,
        code: form.code,
        new_password: form.newPassword
      })
      ElMessage.success('密码重置成功，请登录')
      switchMode('login')
    }
  } catch (e) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
.login-card {
  width: 420px;
  padding: 40px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
}
a { color: #409eff; text-decoration: none; }
a:hover { text-decoration: underline; }
</style>
