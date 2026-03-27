<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useTeamStore } from '@/stores/useTeamStore'
import { organizationService } from '@/services/teamService'
import AppEmptyState from '@/components/common/AppEmptyState.vue'
import { getInitials } from '@/utils/format'
import type { Organization } from '@/types/team'
import { useToast } from '@/composables/useToast'
import { invitationService } from '@/services/invitationService'
import { useConfirm } from '@/composables/useConfirm'

const { t } = useI18n()
const toast = useToast()
const { confirm } = useConfirm()
const teamStore = useTeamStore()
const showInviteModal = ref(false)
const showCreateModal = ref(false)
const inviteEmail = ref('')
const inviteRole = ref('member')
const newTeamName = ref('')
const newTeamOrgId = ref('')
const organizations = ref<Organization[]>([])
const inviteLoading = ref(false)
const inviteError = ref('')
const createLoading = ref(false)

const roleColors: Record<string, string> = { owner: '#ef4444', leader: '#6366f1', member: '#10b981', viewer: '#64748b' }
const roleOptions = [{ value: 'member', title: 'Member' }, { value: 'leader', title: 'Leader' }, { value: 'viewer', title: 'Viewer' }]

async function loadTeams() {
  await teamStore.fetchTeams()
  if (teamStore.teams.length > 0) await teamStore.fetchTeam(teamStore.teams[0].id)
  try {
    organizations.value = await organizationService.list()
  } catch {
    // No orgs yet
  }
}
onMounted(loadTeams)

async function handleCreateTeam() {
  if (!newTeamName.value) return
  createLoading.value = true
  try {
    // If no org selected, create a default one
    let orgId = newTeamOrgId.value
    if (!orgId) {
      if (organizations.value.length > 0) {
        orgId = organizations.value[0].id
      } else {
        const newOrg = await organizationService.create('My Organization')
        organizations.value.push(newOrg)
        orgId = newOrg.id
      }
    }
    await teamStore.createTeam(newTeamName.value, orgId)
    toast.success(t('team.createTeam') + ' ✓')
    showCreateModal.value = false
    newTeamName.value = ''
    newTeamOrgId.value = ''
  } finally {
    createLoading.value = false
  }
}

async function selectTeam(teamId: string) { await teamStore.fetchTeam(teamId) }

async function handleInvite() {
  if (!inviteEmail.value || !teamStore.currentTeam) return
  inviteLoading.value = true
  inviteError.value = ''
  try {
    await invitationService.send(teamStore.currentTeam.id, inviteEmail.value, inviteRole.value)
    toast.success(t('team.sendInvite') + ' ✓')
    showInviteModal.value = false
    inviteEmail.value = ''
  } catch (err: unknown) {
    const axiosErr = err as { response?: { data?: { error?: string } } }
    inviteError.value = axiosErr.response?.data?.error ?? t('team.inviteFailed')
    toast.error(inviteError.value)
  }
  finally { inviteLoading.value = false }
}

async function handleDeleteTeam(teamId: string) {
  const ok = await confirm({ title: t('common.delete'), message: t('team.teams') + ' — ' + t('common.confirm') + '?', confirmText: t('common.delete'), danger: true })
  if (!ok) return
  await teamStore.deleteTeam(teamId)
  toast.success(t('common.delete') + ' ✓')
}

async function handleRemoveMember(userId: string) {
  if (!teamStore.currentTeam) return
  const ok = await confirm({ title: t('common.remove'), message: t('common.confirm') + '?', confirmText: t('common.remove'), danger: true })
  if (!ok) return
  await teamStore.removeMember(teamStore.currentTeam.id, userId)
  toast.success(t('common.remove') + ' ✓')
}

async function handleRoleChange(userId: string, role: string) {
  if (!teamStore.currentTeam) return
  await teamStore.updateMemberRole(teamStore.currentTeam.id, userId, role)
  toast.success(t('team.role') + ' ✓')
}
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 space-y-5">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold">{{ t('team.teams') }}</h1>
      <v-btn color="primary" prepend-icon="mdi-plus" size="small" style="text-transform: none" @click="showCreateModal = true">{{ t('team.newTeam') }}</v-btn>
    </div>

    <div v-if="teamStore.loading && !teamStore.teams.length" class="flex justify-center py-16">
      <v-progress-circular indeterminate color="primary" size="40" />
    </div>

    <AppEmptyState v-else-if="!teamStore.teams.length" icon="mdi-account-group-outline" :title="t('team.noTeamsYet')" :description="t('team.noTeamsDescription')" :action-label="t('team.createTeam')" @action="showCreateModal = true" />

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-4">
      <!-- Team list -->
      <div class="space-y-2">
        <div
          v-for="team in teamStore.teams"
          :key="team.id"
          class="flex items-center gap-3 p-3 rounded-xl border transition-all group"
          :style="{ background: teamStore.currentTeam?.id === team.id ? 'var(--bg-active)' : 'var(--bg-card)', borderColor: teamStore.currentTeam?.id === team.id ? 'rgba(99,102,241,0.3)' : 'var(--border)' }"
        >
          <button class="flex items-center gap-3 flex-1 text-left" @click="selectTeam(team.id)">
            <div class="w-9 h-9 rounded-lg bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
              <v-icon icon="mdi-account-group" size="18" class="text-primary-light" />
            </div>
            <div>
              <p class="text-sm font-medium" :style="{ color: 'var(--text)' }">{{ team.name }}</p>
              <p class="text-[11px]" :style="{ color: 'var(--text-muted)' }">{{ t('team.teams') }}</p>
            </div>
          </button>
          <button
            class="p-1.5 rounded-lg opacity-0 group-hover:opacity-100 hover:bg-error/10 transition-all"
            @click="handleDeleteTeam(team.id)"
          >
            <v-icon icon="mdi-delete-outline" size="16" class="text-error/50" />
          </button>
        </div>
      </div>

      <!-- Team detail -->
      <div class="lg:col-span-2">
        <div v-if="teamStore.currentTeam" class="rounded-xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
          <div class="flex items-center justify-between px-4 py-3 border-b border-white/5">
            <h2 class="text-sm font-semibold">{{ teamStore.currentTeam.name }}</h2>
            <v-btn size="small" variant="tonal" color="primary" prepend-icon="mdi-account-plus-outline" style="text-transform: none" @click="showInviteModal = true">{{ t('team.invite') }}</v-btn>
          </div>

          <div v-if="teamStore.currentTeam.members?.length" class="divide-y divide-white/[0.03]">
            <div v-for="member in teamStore.currentTeam.members" :key="member.id" class="flex items-center justify-between px-4 py-3">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center">
                  <span class="text-[10px] text-white font-semibold">{{ getInitials(member.user?.name ?? '') }}</span>
                </div>
                <div>
                  <p class="text-sm font-medium">{{ member.user?.name }}</p>
                  <p class="text-[11px] text-white/25">{{ member.user?.email }}</p>
                </div>
              </div>
              <div class="flex items-center gap-2">
                <span class="px-2 py-0.5 rounded-full text-[10px] font-medium" :style="{ backgroundColor: (roleColors[member.role] ?? '#64748b') + '18', color: roleColors[member.role] }">{{ member.role }}</span>
                <v-menu>
                  <template #activator="{ props }">
                    <button v-bind="props" class="p-1 rounded hover:bg-white/5 transition-colors">
                      <v-icon icon="mdi-dots-vertical" size="16" class="text-white/30" />
                    </button>
                  </template>
                  <v-list density="compact" min-width="140" bg-color="surface" class="rounded-xl border border-white/5">
                    <v-list-item v-for="role in roleOptions" :key="role.value" :title="t('team.setRole', { role: role.title })" @click="handleRoleChange(member.userId, role.value)" />
                    <v-divider class="opacity-10" />
                    <v-list-item :title="t('common.remove')" class="text-error" @click="handleRemoveMember(member.userId)" />
                  </v-list>
                </v-menu>
              </div>
            </div>
          </div>
          <div v-else class="py-8 text-center text-xs text-white/25">{{ t('team.noMembers') }}</div>
        </div>

        <AppEmptyState v-else icon="mdi-account-group-outline" :title="t('team.selectTeam')" :description="t('team.selectTeamDescription')" />
      </div>
    </div>

    <!-- Invite modal -->
    <v-dialog v-model="showInviteModal" max-width="400">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-3">{{ t('team.inviteMember') }}</h3>
        <v-alert v-if="inviteError" type="error" variant="tonal" class="mb-3 text-xs">{{ inviteError }}</v-alert>
        <v-text-field v-model="inviteEmail" :label="t('team.emailAddress')" type="email" prepend-inner-icon="mdi-email-outline" class="mb-2" />
        <v-select v-model="inviteRole" :items="roleOptions" :label="t('team.role')" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showInviteModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :loading="inviteLoading" style="text-transform: none" @click="handleInvite">{{ t('team.sendInvite') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Create team modal -->
    <v-dialog v-model="showCreateModal" max-width="400">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-3">{{ t('team.createTeam') }}</h3>
        <v-text-field v-model="newTeamName" :label="t('team.teamName')" prepend-inner-icon="mdi-account-group-outline" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showCreateModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :disabled="!newTeamName" :loading="createLoading" style="text-transform: none" @click="handleCreateTeam">{{ t('common.create') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
