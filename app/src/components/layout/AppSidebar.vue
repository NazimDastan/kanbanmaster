<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const route = useRoute()
const collapsed = ref(false)

const navItems = computed(() => [
  { title: t('sidebar.dashboard'), icon: 'mdi-view-dashboard-outline', to: '/' },
  { title: t('sidebar.boards'), icon: 'mdi-view-column-outline', to: '/boards' },
  { title: t('sidebar.teams'), icon: 'mdi-account-group-outline', to: '/teams' },
  { title: t('sidebar.reports'), icon: 'mdi-chart-bar', to: '/reports' },
])

function isActive(path: string): boolean {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}
</script>

<template>
  <aside
    class="flex-shrink-0 flex flex-col border-r transition-all duration-200"
    :style="{ background: 'var(--bg-card)', borderColor: 'var(--border)' }"
    :class="collapsed ? 'w-16' : 'w-56'"
  >
    <nav class="flex-1 overflow-y-auto px-2 py-3 space-y-0.5">
      <router-link
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm transition-all duration-150"
        :style="isActive(item.to)
          ? { background: 'var(--bg-active)', color: '#818cf8' }
          : { color: 'var(--text-secondary)' }"
      >
        <v-icon :icon="item.icon" size="18" />
        <span v-if="!collapsed" class="truncate">{{ item.title }}</span>
      </router-link>
    </nav>

    <div class="p-2" :style="{ borderTop: '1px solid var(--border)' }">
      <button
        class="w-full flex items-center justify-center py-1.5 rounded-lg transition-all"
        :style="{ color: 'var(--text-muted)' }"
        @click="collapsed = !collapsed"
      >
        <v-icon :icon="collapsed ? 'mdi-chevron-right' : 'mdi-chevron-left'" size="18" />
      </button>
    </div>
  </aside>
</template>
