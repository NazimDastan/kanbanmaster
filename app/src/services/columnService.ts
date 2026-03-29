import api from './api'
import type { Column } from '@/types/board'

export const columnService = {
  async create(boardId: string, name: string, color?: string): Promise<Column> {
    const { data } = await api.post<Column>(`/boards/${boardId}/columns`, { name, color: color ?? '' })
    return data
  },

  async update(id: string, name: string): Promise<Column> {
    const { data } = await api.put<Column>(`/columns/${id}`, { name })
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/columns/${id}`)
  },

  async reorder(boardId: string, items: { columnId: string; position: number }[]): Promise<void> {
    await api.patch('/columns/reorder', { boardId, items })
  },
}
