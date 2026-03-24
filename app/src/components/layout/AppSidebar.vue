<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useBoardStore } from '@/stores/useBoardStore'
import { boardService } from '@/services/boardService'
import { organizationService, teamService } from '@/services/teamService'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const boardStore = useBoardStore()
const collapsed = ref(false)
const showNewBoard = ref(false)
const newBoardName = ref('')
const creating = ref(false)

const navItems = computed(() => [
  { title: t('sidebar.dashboard'), icon: 'mdi-view-dashboard-outline', to: '/' },
  { title: t('sidebar.teams'), icon: 'mdi-account-group-outline', to: '/teams' },
  { title: t('sidebar.reports'), icon: 'mdi-chart-bar', to: '/reports' },
])

function isActive(path: string): boolean {
  if (path === '/' ) return route.path === '/'
  return route.path.startsWith(path)
}

async function loadBoards() {
  await boardStore.fetchBoards()
}
onMounted(loadBoards)

async function handleCreateBoard() {
  if (!newBoardName.value.trim()) return
  creating.value = true
  try {
    // Ensure team exists
    let teams = await teamService.list()
    let teamId: string

    if (teams.length > 0) {
      teamId = teams[0].id
    } else {
      // Create org + team automatically
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
    newBoardName.value = ''
    showNewBoard.value = false
    router.push(`/boards/${board.id}`)
  } finally {
    creating.value = false
  }
}
</script>

<template>
  <aside
    class="flex-shrink-0 flex flex-col border-r transition-all duration-200"
    style="background: var(--bg-card, #0f0f1a); border-color: var(--border, rgba(255,255,255,0.05))"
    :class="collapsed ? 'w-16' : 'w-56'"
  >
    <nav class="flex-1 overflow-y-auto px-2 py-3 space-y-0.5">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150"
        :class="isActive(item.to) ? 'bg-primary/15 text-primary-light' : 'text-white/50 hover:text-white/80 hover:bg-white/5'"
      >
        <v-icon :icon="item.icon" size="18" />
        <span v-if="!collapsed" class="truncate">{{ item.title }}</span>
      </router-link>

      <!-- Boards section -->
      <div class="pt-4 mt-4 border-t border-white/5">
        <div v-if="!collapsed" class="flex items-center justify-between px-3 mb-2">
          <span class="text-[10px] font-semibold uppercase tracking-widest text-white/30">{{ t('sidebar.boards') }}</span>
          <button class="text-white/30 hover:text-primary transition-colors" @click="showNewBoard = !showNewBoard">
            <v-icon :icon="showNewBoard ? 'mdi-close' : 'mdi-plus'" size="14" />
          </button>
        </div>

        <!-- Inline board creation -->
        <div v-if="showNewBoard && !collapsed" class="px-2 mb-2">
          <div class="flex gap-1">
            <input
              v-model="newBoardName"
              :placeholder="t('board.columnName')"
              class="flex-1 px-2.5 py-1.5 rounded-lg bg-white/5 border border-white/10 text-xs text-white placeholder:text-white/20 outline-none focus:border-primary/40"
              @keyup.enter="handleCreateBoard"
            />
            <button
              class="px-2 py-1.5 rounded-lg bg-primary/20 text-primary-light hover:bg-primary/30 transition-colors disabled:opacity-30"
              :disabled="!newBoardName.trim() || creating"
              @click="handleCreateBoard"
            >
              <v-icon icon="mdi-check" size="14" />
            </button>
          </div>
        </div>

        <!-- Board list -->
        <button
          v-for="board in boardStore.boards"
          :key="board.id"
          class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all text-left"
          :class="route.path === `/boards/${board.id}` ? 'bg-primary/15 text-primary-light' : 'text-white/50 hover:text-white/80 hover:bg-white/5'"
          @click="router.push(`/boards/${board.id}`)"
        >
          <v-icon icon="mdi-view-column-outline" size="16" />
          <span v-if="!collapsed" class="truncate">{{ board.name }}</span>
        </button>

        <!-- Empty state -->
        <div v-if="boardStore.boards.length === 0 && !collapsed && !showNewBoard" class="px-3 py-3">
          <p class="text-[11px] text-white/20 mb-2">{{ t('sidebar.noBoardsYet') }}</p>
          <button
            class="w-full flex items-center gap-2 px-3 py-2 rounded-lg border border-dashed border-white/10 text-[11px] text-white/30 hover:text-primary-light hover:border-primary/30 transition-all"
            @click="showNewBoard = true"
          >
            <v-icon icon="mdi-plus" size="14" />
            {{ t('sidebar.createFirstBoard') }}
          </button>
        </div>
      </div>
    </nav>

    <div class="p-2 border-t border-white/5">
      <button class="w-full flex items-center justify-center py-1.5 rounded-lg text-white/30 hover:text-white/60 hover:bg-white/5 transition-all" @click="collapsed = !collapsed">
        <v-icon :icon="collapsed ? 'mdi-chevron-right' : 'mdi-chevron-left'" size="18" />
      </button>
    </div>
  </aside>
</template>
