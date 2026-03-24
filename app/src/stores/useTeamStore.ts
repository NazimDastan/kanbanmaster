import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Team, TeamWithMembers } from '@/types/team'
import { teamService } from '@/services/teamService'

export const useTeamStore = defineStore('team', () => {
  const teams = ref<Team[]>([])
  const currentTeam = ref<TeamWithMembers | null>(null)
  const loading = ref(false)

  async function fetchTeams() {
    loading.value = true
    try {
      teams.value = await teamService.list()
    } finally {
      loading.value = false
    }
  }

  async function fetchTeam(id: string) {
    loading.value = true
    try {
      currentTeam.value = await teamService.get(id)
    } finally {
      loading.value = false
    }
  }

  async function createTeam(name: string, organizationId: string) {
    const team = await teamService.create(name, organizationId)
    teams.value.push(team)
    return team
  }

  async function deleteTeam(id: string) {
    await teamService.delete(id)
    teams.value = teams.value.filter((t) => t.id !== id)
    if (currentTeam.value?.id === id) {
      currentTeam.value = null
    }
  }

  async function inviteMember(teamId: string, email: string, role?: string) {
    const member = await teamService.invite(teamId, email, role)
    if (currentTeam.value?.id === teamId) {
      currentTeam.value.members.push(member)
    }
    return member
  }

  async function removeMember(teamId: string, userId: string) {
    await teamService.removeMember(teamId, userId)
    if (currentTeam.value?.id === teamId) {
      currentTeam.value.members = currentTeam.value.members.filter(
        (m) => m.userId !== userId,
      )
    }
  }

  async function updateMemberRole(teamId: string, userId: string, role: string) {
    await teamService.updateMemberRole(teamId, userId, role)
    if (currentTeam.value?.id === teamId) {
      const member = currentTeam.value.members.find((m) => m.userId === userId)
      if (member) member.role = role
    }
  }

  return {
    teams, currentTeam, loading,
    fetchTeams, fetchTeam, createTeam, deleteTeam,
    inviteMember, removeMember, updateMemberRole,
  }
})
