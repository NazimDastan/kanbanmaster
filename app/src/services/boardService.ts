import api from './api'
import type { Board, BoardWithColumns } from '@/types/board'

export const boardService = {
  async list(): Promise<Board[]> {
    const response = await api.get<Board[]>('/boards')
    return response.data
  },

  async get(id: string): Promise<BoardWithColumns> {
    const response = await api.get<BoardWithColumns>(`/boards/${id}`)
    return response.data
  },

  async create(data: { name: string; teamId: string }): Promise<Board> {
    const response = await api.post<Board>('/boards', data)
    return response.data
  },

  async update(id: string, data: { name: string }): Promise<Board> {
    const response = await api.put<Board>(`/boards/${id}`, data)
    return response.data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/boards/${id}`)
  },
}
