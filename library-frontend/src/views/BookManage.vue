<template>
  <div>
    <!-- 工具栏 -->
    <el-card shadow="never" style="margin-bottom:16px;">
      <div style="display:flex;justify-content:space-between;">
        <el-button type="primary" @click="openCreateDialog">
          <el-icon><Plus /></el-icon> 新增图书
        </el-button>
      </div>
    </el-card>

    <!-- 图书列表 -->
    <el-card shadow="never">
      <el-table :data="books" v-loading="loading" stripe style="width:100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="title" label="书名" min-width="160" />
        <el-table-column prop="author" label="作者" width="130" />
        <el-table-column prop="isbn" label="ISBN" width="140" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'warning'" size="small">
              {{ row.status === 1 ? '可借阅' : '已借出' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确定删除吗？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div style="display:flex;justify-content:center;margin-top:16px;">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.size"
          :total="total"
          layout="total, prev, pager, next"
          @change="fetchBooks"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑图书' : '新增图书'" width="500px">
      <el-form :model="dialogForm" label-width="80px">
        <el-form-item label="书名" required>
          <el-input v-model="dialogForm.title" />
        </el-form-item>
        <el-form-item label="作者" required>
          <el-input v-model="dialogForm.author" />
        </el-form-item>
        <el-form-item label="ISBN" required>
          <el-input v-model="dialogForm.isbn" />
        </el-form-item>
        <el-form-item label="出版日期">
          <el-input v-model="dialogForm.publish_date" placeholder="2006-01-02" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { getBooksApi, createBookApi, updateBookApi, deleteBookApi } from '../api/book'
import { ElMessage } from 'element-plus'

const loading = ref(false)
const saving = ref(false)
const books = ref([])
const total = ref(0)
const query = reactive({ page: 1, size: 10 })
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(null)
const dialogForm = reactive({ title: '', author: '', isbn: '', publish_date: '' })

async function fetchBooks() {
  loading.value = true
  try {
    const res = await getBooksApi({ page: query.page, size: query.size })
    books.value = res.data.list || []
    total.value = res.data.total || 0
  } finally {
    loading.value = false
  }
}

function openCreateDialog() {
  isEdit.value = false
  editId.value = null
  dialogForm.title = ''
  dialogForm.author = ''
  dialogForm.isbn = ''
  dialogForm.publish_date = ''
  dialogVisible.value = true
}

function handleEdit(row) {
  isEdit.value = true
  editId.value = row.id
  dialogForm.title = row.title
  dialogForm.author = row.author
  dialogForm.isbn = row.isbn
  dialogForm.publish_date = row.publish_date || ''
  dialogVisible.value = true
}

async function handleSave() {
  if (!dialogForm.title || !dialogForm.author || !dialogForm.isbn) {
    ElMessage.warning('请填写完整信息')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateBookApi(editId.value, dialogForm)
      ElMessage.success('更新成功')
    } else {
      await createBookApi(dialogForm)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchBooks()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await deleteBookApi(id)
  ElMessage.success('删除成功')
  fetchBooks()
}

onMounted(fetchBooks)
</script>
