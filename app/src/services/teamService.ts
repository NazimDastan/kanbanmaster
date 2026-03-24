import api from './api'
import type { Organization } from '@/types/team'
import type { Team, TeamWithMembers, TeamMember } from '@/types/team'

export const organizationService = {
  async list(): Promise<Organization[]> {
    const { data } = await api.get<Organization[]>('/organizations')
    return data
  },

  async get(id: string): Promise<Organization> {
    const { data } = await api.get<Organization>(`/organizations/${id}`)
    return data
  },

  async create(name: string): Promise<Organization> {
    const { data } = await api.post<Organization>('/organizations', { name })
    return data
  },

  async update(id: string, name: string): Promise<Organization> {
    const { data } = await api.put<Organization>(`/organizations/${id}`, { name })
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/organizations/${id}`)
  },
}

export const teamService = {
  async list(): Promise<Team[]> {
    const { data } = await api.get<Team[]>('/teams')
    return data
  },

  async get(id: string): Promise<TeamWithMembers> {
    const { data } = await api.get<TeamWithMembers>(`/teams/${id}`)
    return data
  },

  async create(name: string, organizationId: string): Promise<Team> {
    const { data } = await api.post<Team>('/teams', { name, organizationId })
    return data
  },

  async update(id: string, name: string): Promise<Team> {
    const { data } = await api.put<Team>(`/teams/${id}`, { name })
    return data
  },

  async delete(id: string): Promise<void> {
    await api.delete(`/teams/${id}`)
  },

  async invite(teamId: string, email: string, role?: string): Promise<TeamMember> {
    const { data } = await api.post<TeamMember>(`/teams/${teamId}/invite`, { email, role })
    return data
  },

  async removeMember(teamId: string, userId: string): Promise<void> {
    await api.delete(`/teams/${teamId}/members/${userId}`)
  },

  async updateMemberRole(teamId: string, userId: string, role: string): Promise<void> {
    await api.patch(`/teams/${teamId}/members/${userId}/role`, { role })
  },
}
