import api from './api'
import type { Label } from '@/types/task'

export const labelService = {
  async listByBoard(boardId: string): Promise<Label[]> {
    const { data } = await api.get<Label[]>(`/boards/${boardId}/labels`)
    return data
  },

  async create(boardId: string, name: string, color: string): Promise<Label> {
    const { data } = await api.post<Label>(`/boards/${boardId}/labels`, { name, color })
    return data
  },

  async update(id: string, name: string, color: string): Promise<Label> {
    const { data } = await api.put<Label>(`/labels/${id}`, { name, color })
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/labels/${id}`)
  },

  async addToTask(taskId: string, labelId: string): Promise<void> {
    await api.post(`/tasks/${taskId}/labels`, { labelId })
  },

  async removeFromTask(taskId: string, labelId: string): Promise<void> {
    await api.delete(`/tasks/${taskId}/labels/${labelId}`)
  },
}
