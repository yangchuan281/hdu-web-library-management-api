<template>
  <div>
    <!-- 借阅操作栏 -->
    <el-card shadow="never" style="margin-bottom:16px;">
      <div style="display:flex;gap:12px;align-items:center;">
        <el-input v-model="borrowBookId" placeholder="输入图书ID借阅" style="width:200px;" />
        <el-button type="primary" @click="handleBorrow" :disabled="!borrowBookId">
          <el-icon><Plus /></el-icon> 借阅
        </el-button>
      </div>
    </el-card>

    <!-- 借阅列表 -->
    <el-card shadow="never">
      <el-table :data="borrows" v-loading="loading" stripe style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="book_id" label="图书ID" width="80" />
        <el-table-column prop="book_name" label="书名" min-width="160" />
        <el-table-column prop="borrow_at" label="借阅时间" width="170" />
        <el-table-column prop="return_at" label="归还时间" width="170">
          <template #default="{ row }">
            {{ row.return_at || '借阅中' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="!row.return_at"
              size="small"
              type="success"
              @click="handleReturn(row.id)"
            >归还</el-button>
            <el-tag v-else type="info" size="small">已归还</el-tag>
          </template>
        </el-table-column>
      </el-table>

      <div style="display:flex;justify-content:center;margin-top:16px;">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          layout="total, prev, pager, next"
          @change="fetchBorrows"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBorrowsApi, borrowBookApi, returnBookApi } from '../api/borrow'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const borrows = ref([])
const total = ref(0)
const borrowBookId = ref('')
const query = reactive({ page: 1, size: 10 })

async function fetchBorrows() {
  loading.value = true
  try {
    const res = await getBorrowsApi({ page: query.page, size: query.size })
    borrows.value = res.data.list || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

async function handleBorrow() {
  if (!borrowBookId.value) return
  try {
    await borrowBookApi(Number(borrowBookId.value))
    ElMessage.success('借阅成功')
    borrowBookId.value = ''
    fetchBorrows()
  } catch (e) {
    // 错误已在拦截器中处理
  }
}

async function handleReturn(id) {
  try {
    await returnBookApi(id)
    ElMessage.success('归还成功')
    fetchBorrows()
  } catch (e) {}
}

onMounted(fetchBorrows)
</script>
