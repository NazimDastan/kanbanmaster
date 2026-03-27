<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { invitationService, type Invitation } from '@/services/invitationService'
import { useToast } from '@/composables/useToast'
import { formatRelativeTime } from '@/utils/date'

const { t } = useI18n()
const toast = useToast()

const invitations = ref<Invitation[]>([])
const loading = ref(false)

async function loadInvitations() {
  loading.value = true
  try { invitations.value = await invitationService.getPending() }
  catch { /* */ }
  finally { loading.value = false }
}
onMounted(loadInvitations)

async function handleAccept(id: string) {
  try {
    await invitationService.accept(id)
    invitations.value = invitations.value.filter(i => i.id !== id)
    toast.success(t('common.confirm') + ' ✓')
  } catch {
    toast.error('Failed')
  }
}

async function handleReject(id: string) {
  try {
    await invitationService.reject(id)
    invitations.value = invitations.value.filter(i => i.id !== id)
    toast.info(t('common.cancel'))
  } catch {
    toast.error('Failed')
  }
}

defineExpose({ invitations, loadInvitations })
</script>

<template>
  <v-menu :close-on-content-click="false" max-width="380">
    <template #activator="{ props }">
      <button v-bind="props" class="relative p-2 rounded-lg transition-colors">
        <v-icon icon="mdi-email-outline" size="20" :style="{ color: 'var(--text-secondary)' }" />
        <span
          v-if="invitations.length > 0"
          class="absolute top-1 right-1 w-4 h-4 rounded-full bg-warning text-[9px] font-bold text-white flex items-center justify-center"
        >
          {{ invitations.length }}
        </span>
      </button>
    </template>

    <v-card min-width="350" max-height="400" color="surface" class="rounded-xl">
      <div class="px-4 py-3" :style="{ borderBottom: '1px solid var(--border)' }">
        <span class="text-xs font-semibold" :style="{ color: 'var(--text)' }">{{ t('team.inviteMember') }}</span>
      </div>

      <div v-if="invitations.length > 0" class="overflow-y-auto" style="max-height: 320px">
        <div
          v-for="inv in invitations"
          :key="inv.id"
          class="px-4 py-3"
          :style="{ borderBottom: '1px solid var(--border)' }"
        >
          <div class="flex items-center gap-3 mb-2">
            <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
              <v-icon icon="mdi-account-group" size="16" color="white" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium" :style="{ color: 'var(--text)' }">{{ inv.team?.name }}</p>
              <p class="text-[11px]" :style="{ color: 'var(--text-secondary)' }">
                {{ inv.inviter?.name }} · {{ formatRelativeTime(inv.createdAt) }}
              </p>
            </div>
          </div>
          <div class="flex gap-2">
            <v-btn size="small" color="primary" variant="flat" class="flex-1" @click="handleAccept(inv.id)">
              {{ t('common.confirm') }}
            </v-btn>
            <v-btn size="small" variant="outlined" class="flex-1" @click="handleReject(inv.id)">
              {{ t('common.cancel') }}
            </v-btn>
          </div>
        </div>
      </div>

      <div v-else class="py-8 text-center">
        <v-icon icon="mdi-email-check-outline" size="32" :style="{ color: 'var(--text-faint)' }" />
        <p class="text-xs mt-2" :style="{ color: 'var(--text-muted)' }">{{ t('common.noData') }}</p>
      </div>
    </v-card>
  </v-menu>
</template>
