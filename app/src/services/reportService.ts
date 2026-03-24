import api from './api'

export interface ReportRequest {
  id: string
  requesterId: string
  requester?: { id: string; name: string; email: string }
  targetUserId: string
  targetUser?: { id: string; name: string; email: string }
  teamId: string
  message: string
  response: string | null
  status: 'pending' | 'submitted' | 'reviewed'
  createdAt: string
  respondedAt: string | null
}

export interface ActivityLog {
  id: string
  taskId: string
  userId: string
  user?: { id: string; name: string; email: string }
  action: string
  details: Record<string, unknown>
  createdAt: string
}

export const reportService = {
  async requestReport(data: { targetUserId: string; teamId: string; message: string }): Promise<ReportRequest> {
    const { data: result } = await api.post<ReportRequest>('/reports/request', data)
    return result
  },

  async getIncoming(): Promise<ReportRequest[]> {
    const { data } = await api.get<ReportRequest[]>('/reports/requests')
    return data
  },

  async getSent(): Promise<ReportRequest[]> {
    const { data } = await api.get<ReportRequest[]>('/reports/requests/sent')
    return data
  },

  async respond(id: string, response: string): Promise<ReportRequest> {
    const { data } = await api.post<ReportRequest>(`/reports/requests/${id}/respond`, { response })
    return data
  },

  async review(id: string): Promise<ReportRequest> {
    const { data } = await api.patch<ReportRequest>(`/reports/requests/${id}/review`)
    return data
  },
}

export const delegationService = {
  async delegate(taskId: string, toUserId: string, reason: string): Promise<unknown> {
    const { data } = await api.post(`/tasks/${taskId}/delegate`, { toUserId, reason })
    return data
  },

  async getActivity(taskId: string): Promise<ActivityLog[]> {
    const { data } = await api.get<ActivityLog[]>(`/tasks/${taskId}/activity`)
    return data
  },
}
