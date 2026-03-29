<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { teamService } from '@/services/teamService'
import { getInitials } from '@/utils/format'
import type { TeamMember } from '@/types/team'

const { t } = useI18n()
const props = defineProps<{ teamId: string; currentAssigneeId?: string | null }>()
const emit = defineEmits<{ select: [userId: string, userName: string] }>()

const members = ref<TeamMember[]>([])
const searchQuery = ref('')
const loading = ref(false)
const showDropdown = ref(false)

const filtered = computed(() => {
  const q = searchQuery.value.toLowerCase().trim()
  if (!q) return members.value
  return members.value.filter((m) =>
    m.user.name.toLowerCase().includes(q) || m.user.email.toLowerCase().includes(q),
  )
})

async function loadMembers() {
  if (!props.teamId) return
  loading.value = true
  try {
    const team = await teamService.get(props.teamId)
    members.value = team.members ?? []
  } catch {
    members.value = []
  } finally {
    loading.value = false
  }
}

onMounted(loadMembers)
watch(() => props.teamId, loadMembers)

function handleSelect(member: TeamMember) {
  emit('select', member.userId, member.user.name)
  showDropdown.value = false
  searchQuery.value = ''
}
</script>

<template>
  <div class="relative">
    <!-- Trigger -->
    <button
      class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-xl border border-dashed transition-all text-sm border-[var(--border)] text-[var(--text-muted)] hover:border-primary/30 hover:bg-primary/5"
      @click="showDropdown = !showDropdown"
    >
      <v-icon icon="mdi-account-plus-outline" size="18" />
      {{ t('task.assignTask') }}
    </button>

    <!-- Dropdown -->
    <div
      v-if="showDropdown"
      class="absolute top-full left-0 right-0 mt-1.5 z-50 rounded-xl border bg-[var(--bg-card)] border-[var(--border)] shadow-xl overflow-hidden"
    >
      <!-- Search -->
      <div class="p-2 border-b border-[var(--border)]">
        <div class="relative">
          <v-icon icon="mdi-magnify" size="15" class="absolute left-2.5 top-1/2 -translate-y-1/2 text-[var(--text-muted)]" />
          <input
            v-model="searchQuery"
            :placeholder="t('task.searchMembers')"
            class="w-full pl-8 pr-3 py-2 rounded-lg bg-[var(--bg-input)] text-xs outline-none placeholder:text-[var(--text-placeholder)] text-[var(--text)]"
            autofocus
          />
        </div>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="p-4 text-center">
        <v-progress-circular indeterminate size="20" width="2" color="primary" />
      </div>

      <!-- Members list -->
      <div v-else class="max-h-[240px] overflow-y-auto">
        <div v-if="filtered.length === 0" class="p-4 text-center text-xs text-[var(--text-muted)]">
          {{ t('common.noResults') }}
        </div>
        <button
          v-for="member in filtered"
          :key="member.userId"
          class="w-full flex items-center gap-3 px-3 py-2.5 hover:bg-[var(--bg-hover)] transition-colors text-left"
          :class="{ 'bg-primary/5': member.userId === currentAssigneeId }"
          @click="handleSelect(member)"
        >
          <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
            <span class="text-[9px] text-white font-semibold">{{ getInitials(member.user.name) }}</span>
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-xs font-medium text-[var(--text)] truncate">{{ member.user.name }}</p>
            <p class="text-[10px] text-[var(--text-muted)] truncate">{{ member.user.email }}</p>
          </div>
          <span class="text-[9px] px-1.5 py-0.5 rounded bg-[var(--bg-input)] text-[var(--text-muted)] flex-shrink-0">{{ member.role }}</span>
          <v-icon v-if="member.userId === currentAssigneeId" icon="mdi-check" size="14" class="text-primary flex-shrink-0" />
        </button>
      </div>
    </div>

    <!-- Click outside to close -->
    <div v-if="showDropdown" class="fixed inset-0 z-40" @click="showDropdown = false" />
  </div>
</template>
