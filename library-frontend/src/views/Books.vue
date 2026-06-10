<template>
  <div>
    <!-- 搜索栏 -->
    <el-card shadow="never" style="margin-bottom:16px;">
      <el-form :inline="true" :model="query" size="default">
        <el-form-item label="书名">
          <el-input v-model="query.title" placeholder="搜索书名" clearable @clear="search" @keyup.enter="search" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="query.status" placeholder="全部" clearable @change="search" style="width:120px;">
            <el-option label="可借阅" :value="1" />
            <el-option label="已借出" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="search">搜索</el-button>
          <el-button @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 图书列表 -->
    <el-card shadow="never">
      <el-table :data="books" v-loading="loading" stripe style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="书名" min-width="180" />
        <el-table-column prop="author" label="作者" width="150" />
        <el-table-column prop="isbn" label="ISBN" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
              {{ row.status === 1 ? '可借阅' : '已借出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="入库时间" width="170" />
      </el-table>

      <!-- 分页 -->
      <div style="display:flex;justify-content:center;margin-top:16px;">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @change="fetchBooks"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBooksApi } from '../api/book'

const loading = ref(false)
const books = ref([])
const total = ref(0)
const query = reactive({
  page: 1,
  size: 10,
  title: '',
  status: undefined
})

async function fetchBooks() {
  loading.value = true
  try {
    const params = { page: query.page, size: query.size }
    if (query.title) params.title = query.title
    if (query.status !== undefined && query.status !== '') params.status = query.status
    const res = await getBooksApi(params)
    books.value = res.data.list || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

function search() {
  query.page = 1
  fetchBooks()
}

function resetQuery() {
  query.title = ''
  query.status = undefined
  query.page = 1
  fetchBooks()
}

onMounted(fetchBooks)
</script>
