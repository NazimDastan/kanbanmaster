import api from './api'

export interface Attachment {
  id: string
  taskId: string
  userId: string
  filename: string
  contentType: string
  size: number
  data: string
  createdAt: string
}

export const attachmentService = {
  async list(taskId: string): Promise<Attachment[]> {
    const { data } = await api.get<Attachment[]>(`/tasks/${taskId}/attachments`)
    return data
  },

  async upload(taskId: string, file: { filename: string; contentType: string; size: number; data: string }): Promise<Attachment> {
    const { data } = await api.post<Attachment>(`/tasks/${taskId}/attachments`, file)
    return data
  },

  async download(id: string): Promise<Attachment> {
    const { data } = await api.get<Attachment>(`/attachments/${id}`)
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/attachments/${id}`)
  },
}
