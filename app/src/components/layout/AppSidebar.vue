<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useBoardStore } from '@/stores/useBoardStore'
import { useAuthStore } from '@/stores/useAuthStore'
import { useConfirm } from '@/composables/useConfirm'
import { getInitials } from '@/utils/format'

const { t } = useI18n()
const route = useRoute()
const boardStore = useBoardStore()
const authStore = useAuthStore()
const { confirm } = useConfirm()
const emit = defineEmits<{ navigate: [] }>()
const collapsed = ref(false)

const navItems = computed(() => [
  { title: t('sidebar.dashboard'), icon: 'mdi-view-dashboard-outline', to: '/' },
  { title: t('sidebar.boards'), icon: 'mdi-view-column-outline', to: '/boards' },
  { title: t('sidebar.teams'), icon: 'mdi-account-group-outline', to: '/teams' },
  { title: t('sidebar.reports'), icon: 'mdi-chart-bar', to: '/reports' },
])

const recentBoards = computed(() => boardStore.boards.slice(0, 5))
const initials = computed(() => getInitials(authStore.user?.name ?? ''))

function isActive(path: string): boolean {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

async function handleLogout() {
  const ok = await confirm({
    title: t('navbar.logout'),
    message: t('navbar.logoutConfirm'),
    confirmText: t('navbar.logout'),
    danger: true,
  })
  if (ok) authStore.logout()
}

onMounted(() => {
  boardStore.fetchBoards()
})
</script>

<template>
  <aside
    class="h-full flex-shrink-0 flex flex-col border-r transition-all duration-200 bg-[var(--bg-card)] border-[var(--border)]"
    :class="collapsed ? 'w-16' : 'w-56'"
  >
    <!-- Navigation -->
    <nav class="flex-1 overflow-y-auto px-2 py-3 space-y-0.5">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150"
        :style="isActive(item.to)
          ? { background: 'var(--bg-active)', color: '#818cf8' }
          : { color: 'var(--text-secondary)' }"
        @click="emit('navigate')"
      >
        <v-icon :icon="item.icon" size="18" />
        <span v-if="!collapsed" class="truncate">{{ item.title }}</span>
      </router-link>

      <!-- Recent Boards -->
      <div class="mt-4 pt-3 border-t border-[var(--border)]">
        <p v-if="!collapsed" class="px-3 mb-1 text-[10px] font-semibold uppercase tracking-widest text-[var(--text-muted)]">
          {{ t('sidebar.recentBoards') }}
        </p>

        <template v-if="recentBoards.length > 0">
          <router-link
            v-for="board in recentBoards"
            :key="board.id"
            :to="`/boards/${board.id}`"
            class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150"
            :style="route.path === `/boards/${board.id}`
              ? { background: 'var(--bg-active)', color: '#818cf8' }
              : { color: 'var(--text-secondary)' }"
            @click="emit('navigate')"
          >
            <v-icon icon="mdi-view-column-outline" size="16" />
            <span v-if="!collapsed" class="truncate">{{ board.name }}</span>
          </router-link>
        </template>

        <p v-else-if="!collapsed" class="px-3 py-2 text-xs text-[var(--text-muted)]">
          {{ t('sidebar.noBoardsYet') }}
        </p>
      </div>
    </nav>

    <!-- Bottom: Profile + Collapse + Logout -->
    <div class="border-t border-[var(--border)] p-2 space-y-1">
      <!-- Profile link -->
      <router-link
        to="/profile"
        class="flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150"
        :style="isActive('/profile')
          ? { background: 'var(--bg-active)', color: '#818cf8' }
          : { color: 'var(--text-secondary)' }"
        @click="emit('navigate')"
      >
        <div class="w-5 h-5 rounded-full bg-gradient-to-br from-primary/70 to-secondary/70 flex items-center justify-center flex-shrink-0">
          <span class="text-[7px] text-white font-semibold">{{ initials }}</span>
        </div>
        <span v-if="!collapsed" class="truncate">{{ authStore.userName }}</span>
      </router-link>

      <!-- Logout -->
      <button
        class="w-full flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-all duration-150 text-error/60 hover:bg-error/5 hover:text-error"
        @click="handleLogout"
      >
        <v-icon icon="mdi-logout" size="18" />
        <span v-if="!collapsed">{{ t('navbar.logout') }}</span>
      </button>

      <!-- Collapse toggle -->
      <button
        class="w-full flex items-center justify-center py-1.5 rounded-lg transition-all text-[var(--text-muted)]"
        aria-label="Toggle sidebar" @click="collapsed = !collapsed"
      >
        <v-icon :icon="collapsed ? 'mdi-chevron-right' : 'mdi-chevron-left'" size="18" />
      </button>
    </div>
  </aside>
</template>
