<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useBoardStore } from '@/stores/useBoardStore'
import { boardService } from '@/services/boardService'
import { organizationService, teamService } from '@/services/teamService'
import { useToast } from '@/composables/useToast'
import { useConfirm } from '@/composables/useConfirm'
import AppEmptyState from '@/components/common/AppEmptyState.vue'
import { useTeamStore } from '@/stores/useTeamStore'
import type { Team } from '@/types/team'

const { t } = useI18n()
const router = useRouter()
const boardStore = useBoardStore()
const teamStore = useTeamStore()
const { confirm } = useConfirm()
const toast = useToast()

const showCreate = ref(false)
const newBoardName = ref('')
const selectedTeamId = ref('')
const creating = ref(false)

async function loadBoards() {
  await boardStore.fetchBoards()
  await teamStore.fetchTeams()
  // Default to personal
  selectedTeamId.value = 'personal'
}
onMounted(loadBoards)

async function getOrCreatePersonalTeam(): Promise<string> {
  // Find or create a "Personal" team
  const existing = teamStore.teams.find(t => t.name === 'Personal')
  if (existing) return existing.id

  let orgs = await organizationService.list()
  let orgId: string
  if (orgs.length > 0) {
    orgId = orgs[0].id
  } else {
    const newOrg = await organizationService.create('My Workspace')
    orgId = newOrg.id
  }
  const newTeam = await teamService.create('Personal', orgId)
  await teamStore.fetchTeams()
  return newTeam.id
}

async function handleCreate() {
  if (!newBoardName.value.trim()) return
  creating.value = true
  try {
    let teamId: string

    if (selectedTeamId.value === 'personal' || !selectedTeamId.value) {
      teamId = await getOrCreatePersonalTeam()
    } else {
      teamId = selectedTeamId.value
    }

    const board = await boardService.create({ name: newBoardName.value, teamId })
    boardStore.boards.push(board)
    toast.success(t('common.create') + ' ✓')
    newBoardName.value = ''
    showCreate.value = false
    router.push(`/boards/${board.id}`)
  } finally {
    creating.value = false
  }
}

async function handleDelete(boardId: string) {
  const ok = await confirm({ title: t('common.delete'), message: t('common.confirm') + '?', confirmText: t('common.delete'), danger: true })
  if (!ok) return
  await boardStore.deleteBoard(boardId)
  toast.success(t('common.delete') + ' ✓')
}
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 space-y-5">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold">{{ t('sidebar.boards') }}</h1>
      <v-btn color="primary" prepend-icon="mdi-plus" size="small" @click="showCreate = true">
        {{ t('dashboard.newBoard') }}
      </v-btn>
    </div>

    <!-- Loading -->
    <div v-if="boardStore.loading" class="flex justify-center py-16">
      <v-progress-circular indeterminate color="primary" size="40" />
    </div>

    <!-- Empty -->
    <AppEmptyState
      v-else-if="boardStore.boards.length === 0"
      icon="mdi-view-column-outline"
      :title="t('dashboard.noBoardsYet')"
      :description="t('dashboard.createBoardHint')"
      :action-label="t('dashboard.newBoard')"
      @action="showCreate = true"
    />

    <!-- Board grid -->
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <button
        v-for="board in boardStore.boards"
        :key="board.id"
        class="flex items-center gap-4 p-4 rounded-xl border text-left transition-all group bg-[var(--bg-card)] border-[var(--border)]"
        @click="router.push(`/boards/${board.id}`)"
      >
        <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center flex-shrink-0">
          <v-icon icon="mdi-view-column-outline" size="24" class="text-primary-light" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold truncate text-[var(--text)]">{{ board.name }}</p>
          <div class="flex items-center gap-1.5 mt-0.5">
            <v-icon :icon="teamStore.teams.find(t => t.id === board.teamId)?.name === 'Personal' ? 'mdi-account-outline' : 'mdi-account-group'" size="12" class="text-[var(--text-muted)]" />
            <span class="text-[11px] text-[var(--text-muted)]">{{ teamStore.teams.find(t => t.id === board.teamId)?.name ?? 'Personal' }}</span>
            <span class="text-[10px] text-[var(--text-faint)]">· {{ board.createdAt.slice(0, 10) }}</span>
          </div>
        </div>
        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <button class="p-1.5 rounded-lg hover:bg-error/10 transition-colors" @click.stop="handleDelete(board.id)">
            <v-icon icon="mdi-delete-outline" size="16" class="text-error/50" />
          </button>
          <v-icon icon="mdi-chevron-right" size="18" class="text-[var(--text-faint)]" />
        </div>
      </button>
    </div>

    <!-- Create dialog -->
    <v-dialog v-model="showCreate" max-width="420">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-4">{{ t('dashboard.newBoard') }}</h3>
        <v-text-field
          v-model="newBoardName"
          :label="t('board.columnName')"
          prepend-inner-icon="mdi-view-column-outline"
          autofocus
          class="mb-2"
        />
        <!-- Workspace selector -->
        <div class="mb-3">
          <p class="text-[10px] font-semibold uppercase tracking-widest mb-2 text-[var(--text-muted)]">{{ t('team.teams') }}</p>
          <div class="flex flex-wrap gap-2">
            <!-- Personal option (default) -->
            <button
              class="flex items-center gap-2 px-3 py-2 rounded-lg border text-xs font-medium transition-all"
              :style="{
                borderColor: selectedTeamId === 'personal' ? 'rgba(16,185,129,0.4)' : 'var(--border)',
                background: selectedTeamId === 'personal' ? 'rgba(16,185,129,0.1)' : 'transparent',
                color: selectedTeamId === 'personal' ? '#10b981' : 'var(--text-secondary)',
              }"
              @click="selectedTeamId = 'personal'"
            >
              <v-icon icon="mdi-account-outline" size="14" />
              Personal
            </button>
            <!-- Team options -->
            <button
              v-for="team in teamStore.teams.filter(t => t.name !== 'Personal')"
              :key="team.id"
              class="flex items-center gap-2 px-3 py-2 rounded-lg border text-xs font-medium transition-all"
              :style="{
                borderColor: selectedTeamId === team.id ? 'rgba(99,102,241,0.4)' : 'var(--border)',
                background: selectedTeamId === team.id ? 'rgba(99,102,241,0.1)' : 'transparent',
                color: selectedTeamId === team.id ? '#818cf8' : 'var(--text-secondary)',
              }"
              @click="selectedTeamId = team.id"
            >
              <v-icon icon="mdi-account-group" size="14" />
              {{ team.name }}
            </button>
          </div>
        </div>
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showCreate = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :disabled="!newBoardName.trim()" :loading="creating" @click="handleCreate">{{ t('common.create') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
