import api from './api'
import type { Task, Comment, ActivityLog } from '@/types/task'

export type TaskFilter = 'assigned' | 'completed' | 'overdue' | 'in_progress' | ''

export const taskService = {
  async list(filter?: TaskFilter): Promise<Task[]> {
    const params = filter ? `?filter=${filter}` : ''
    const { data } = await api.get<Task[]>(`/tasks${params}`)
    return data
  },

  async create(data: {
    columnId: string
    title: string
    description?: string
    assigneeId?: string
    priority: string
    deadline?: string
  }): Promise<Task> {
    const response = await api.post<Task>('/tasks', data)
    return response.data
  },

  async get(id: string): Promise<Task> {
    const response = await api.get<Task>(`/tasks/${id}`)
    return response.data
  },

  async update(id: string, data: Partial<Task>): Promise<Task> {
    const response = await api.put<Task>(`/tasks/${id}`, data)
    return response.data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/tasks/${id}`)
  },

  async move(id: string, data: { columnId: string; position: number }): Promise<void> {
    await api.patch(`/tasks/${id}/move`, data)
  },

  async assign(id: string, assigneeId: string): Promise<void> {
    await api.patch(`/tasks/${id}/assign`, { assigneeId })
  },

  async delegate(id: string, data: { toUserId: string; reason: string }): Promise<void> {
    await api.post(`/tasks/${id}/delegate`, data)
  },

  async getActivity(id: string): Promise<ActivityLog[]> {
    const response = await api.get<ActivityLog[]>(`/tasks/${id}/activity`)
    return response.data
  },

  async getComments(taskId: string): Promise<Comment[]> {
    const response = await api.get<Comment[]>(`/tasks/${taskId}/comments`)
    return response.data
  },

  async addComment(taskId: string, content: string): Promise<Comment> {
    const response = await api.post<Comment>(`/tasks/${taskId}/comments`, { content })
    return response.data
  },

  async addAssignee(taskId: string, userId: string): Promise<void> {
    await api.post(`/tasks/${taskId}/assignees`, { userId })
  },

  async removeAssignee(taskId: string, userId: string): Promise<void> {
    await api.delete(`/tasks/${taskId}/assignees/${userId}`)
  },

  async search(query: string): Promise<Task[]> {
    const { data } = await api.get<Task[]>(`/tasks/search?q=${encodeURIComponent(query)}`)
    return data
  },
}
