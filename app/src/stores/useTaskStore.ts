import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Task } from '@/types/task'
import { taskService } from '@/services/taskService'

export const useTaskStore = defineStore('task', () => {
  const selectedTask = ref<Task | null>(null)
  const loading = ref(false)

  async function fetchTask(id: string) {
    loading.value = true
    try {
      selectedTask.value = await taskService.get(id)
    } finally {
      loading.value = false
    }
  }

  async function createTask(data: {
    columnId: string
    title: string
    description?: string
    assigneeId?: string
    priority: string
    deadline?: string
  }) {
    return await taskService.create(data)
  }

  async function updateTask(id: string, data: Partial<Task>) {
    const updated = await taskService.update(id, data)
    if (selectedTask.value?.id === id) {
      selectedTask.value = updated
    }
    return updated
  }

  async function deleteTask(id: string) {
    await taskService.delete(id)
    if (selectedTask.value?.id === id) {
      selectedTask.value = null
    }
  }

  async function moveTask(id: string, columnId: string, position: number) {
    await taskService.move(id, { columnId, position })
  }

  async function assignTask(id: string, assigneeId: string) {
    await taskService.assign(id, assigneeId)
  }

  return {
    selectedTask, loading,
    fetchTask, createTask, updateTask, deleteTask, moveTask, assignTask,
  }
})
