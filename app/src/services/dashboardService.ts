import api from './api'

export interface DashboardSummary {
  totalTasks: number
  completed: number
  inProgress: number
  overdue: number
  completedPct: number
}

export interface UserPerformance {
  userId: string
  userName: string
  totalTasks: number
  completed: number
  onTime: number
  overdue: number
  score: number
}

export const dashboardService = {
  async getSummary(): Promise<DashboardSummary> {
    const { data } = await api.get<DashboardSummary>('/dashboard/summary')
    return data
  },

  async getTeamPerformance(teamId: string): Promise<UserPerformance[]> {
    const { data } = await api.get<UserPerformance[]>(`/dashboard/team/${teamId}/performance`)
    return data
  },
}
