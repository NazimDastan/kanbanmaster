<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { required } from '@/utils/validation'

const { t } = useI18n()

defineProps<{ taskTitle: string; loading?: boolean }>()
const emit = defineEmits<{ submit: [data: { toUserId: string; reason: string }]; cancel: [] }>()

const form = ref(false)
const toUserId = ref('')
const reason = ref('')

function handleSubmit() {
  if (!form.value) return
  emit('submit', { toUserId: toUserId.value, reason: reason.value })
}
</script>

<template>
  <v-card class="pa-6" max-width="480">
    <h3 class="text-lg font-bold mb-1">{{ t('task.delegateTask') }}</h3>
    <p class="text-sm text-text-secondary mb-4">{{ taskTitle }}</p>
    <v-form v-model="form" @submit.prevent="handleSubmit">
      <v-text-field v-model="toUserId" :label="t('task.delegateTo')" :rules="[required('User ID')]" prepend-inner-icon="mdi-account-arrow-right" class="mb-2" />
      <v-textarea v-model="reason" :label="t('task.delegateReason')" :rules="[required(t('task.delegateReason'))]" rows="3" auto-grow variant="outlined" rounded="lg" class="mb-4" />
      <div class="flex justify-end gap-2">
        <v-btn variant="text" @click="emit('cancel')">{{ t('common.cancel') }}</v-btn>
        <v-btn type="submit" color="secondary" :disabled="!form" :loading="loading" prepend-icon="mdi-swap-horizontal">{{ t('task.delegate') }}</v-btn>
      </div>
    </v-form>
  </v-card>
</template>
