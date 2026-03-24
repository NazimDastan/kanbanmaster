import api from './api'
import type { Comment } from '@/types/task'

export const commentService = {
  async listByTask(taskId: string): Promise<Comment[]> {
    const { data } = await api.get<Comment[]>(`/tasks/${taskId}/comments`)
    return data
  },

  async create(taskId: string, content: string): Promise<Comment> {
    const { data } = await api.post<Comment>(`/tasks/${taskId}/comments`, { content })
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/comments/${id}`)
  },
}
