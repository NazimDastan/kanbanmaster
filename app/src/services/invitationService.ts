import api from './api'

export interface Invitation {
  id: string
  teamId: string
  team?: { name: string }
  inviterId: string
  inviter?: { name: string; email: string }
  inviteeId: string
  invitee?: { name: string; email: string }
  role: string
  status: 'pending' | 'accepted' | 'rejected'
  createdAt: string
  respondedAt: string | null
}

export const invitationService = {
  async send(teamId: string, email: string, role?: string): Promise<Invitation> {
    const { data } = await api.post<Invitation>(`/teams/${teamId}/invite`, { email, role })
    return data
  },

  async getPending(): Promise<Invitation[]> {
    const { data } = await api.get<Invitation[]>('/invitations')
    return data
  },

  async accept(id: string): Promise<void> {
    await api.post(`/invitations/${id}/accept`)
  },

  async reject(id: string): Promise<void> {
    await api.post(`/invitations/${id}/reject`)
  },

  async getTeamInvitations(teamId: string): Promise<Invitation[]> {
    const { data } = await api.get<Invitation[]>(`/teams/${teamId}/invitations`)
    return data
  },
}
