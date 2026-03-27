<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useBoardStore } from '@/stores/useBoardStore'
import { boardService } from '@/services/boardService'
import { organizationService, teamService } from '@/services/teamService'
import { useToast } from '@/composables/useToast'
import AppEmptyState from '@/components/common/AppEmptyState.vue'

const { t } = useI18n()
const router = useRouter()
const boardStore = useBoardStore()
const toast = useToast()

const showCreate = ref(false)
const newBoardName = ref('')
const creating = ref(false)

async function loadBoards() {
  await boardStore.fetchBoards()
}
onMounted(loadBoards)

async function handleCreate() {
  if (!newBoardName.value.trim()) return
  creating.value = true
  try {
    let teams = await teamService.list()
    let teamId: string
    if (teams.length > 0) {
      teamId = teams[0].id
    } else {
      let orgs = await organizationService.list()
      let orgId: string
      if (orgs.length > 0) {
        orgId = orgs[0].id
      } else {
        const newOrg = await organizationService.create('My Organization')
        orgId = newOrg.id
      }
      const newTeam = await teamService.create('My Team', orgId)
      teamId = newTeam.id
    }
    const board = await boardService.create({ name: newBoardName.value, teamId })
    boardStore.boards.push(board)
    toast.success(t('board.addColumn') + ' ✓')
    newBoardName.value = ''
    showCreate.value = false
    router.push(`/boards/${board.id}`)
  } finally {
    creating.value = false
  }
}

async function handleDelete(boardId: string) {
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
        class="flex items-center gap-4 p-4 rounded-xl border text-left transition-all group"
        :style="{ background: 'var(--bg-card)', borderColor: 'var(--border)' }"
        @click="router.push(`/boards/${board.id}`)"
      >
        <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center flex-shrink-0">
          <v-icon icon="mdi-view-column-outline" size="24" class="text-primary-light" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-semibold truncate" :style="{ color: 'var(--text)' }">{{ board.name }}</p>
          <p class="text-[11px] mt-0.5" :style="{ color: 'var(--text-muted)' }">{{ board.createdAt.slice(0, 10) }}</p>
        </div>
        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <button class="p-1.5 rounded-lg hover:bg-error/10 transition-colors" @click.stop="handleDelete(board.id)">
            <v-icon icon="mdi-delete-outline" size="16" class="text-error/50" />
          </button>
          <v-icon icon="mdi-chevron-right" size="18" :style="{ color: 'var(--text-faint)' }" />
        </div>
      </button>
    </div>

    <!-- Create dialog -->
    <v-dialog v-model="showCreate" max-width="400">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-3">{{ t('dashboard.newBoard') }}</h3>
        <v-text-field
          v-model="newBoardName"
          :label="t('board.columnName')"
          prepend-inner-icon="mdi-view-column-outline"
          autofocus
          @keyup.enter="handleCreate"
        />
        <div class="flex justify-end gap-2 mt-3">
          <v-btn variant="text" size="small" @click="showCreate = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :disabled="!newBoardName.trim()" :loading="creating" @click="handleCreate">{{ t('common.create') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
