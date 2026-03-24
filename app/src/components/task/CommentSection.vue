<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Comment } from '@/types/task'
import { commentService } from '@/services/commentService'
import { formatRelativeTime } from '@/utils/date'
import { getInitials } from '@/utils/format'

const { t } = useI18n()
const props = defineProps<{ taskId: string }>()

const comments = ref<Comment[]>([])
const newComment = ref('')
const loading = ref(false)
const sending = ref(false)

async function loadComments() {
  loading.value = true
  try { comments.value = await commentService.listByTask(props.taskId) }
  catch { /* Not connected */ }
  finally { loading.value = false }
}
onMounted(loadComments)

async function handleSubmit() {
  if (!newComment.value.trim()) return
  sending.value = true
  try { const c = await commentService.create(props.taskId, newComment.value); comments.value.push(c); newComment.value = '' }
  finally { sending.value = false }
}

async function handleDelete(commentId: string) {
  await commentService.delete(commentId)
  comments.value = comments.value.filter((c) => c.id !== commentId)
}
</script>

<template>
  <div>
    <p class="text-xs font-semibold uppercase tracking-wider text-text-secondary mb-3">{{ t('task.comments') }} ({{ comments.length }})</p>
    <div class="space-y-3 mb-4">
      <div v-for="comment in comments" :key="comment.id" class="flex gap-3">
        <v-avatar color="primary" size="28" class="mt-0.5">
          <span class="text-[10px] text-white font-medium">{{ getInitials(comment.user?.name ?? 'U') }}</span>
        </v-avatar>
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium">{{ comment.user?.name ?? 'Unknown' }}</span>
            <span class="text-xs text-text-secondary">{{ formatRelativeTime(comment.createdAt) }}</span>
            <v-btn icon="mdi-delete-outline" size="x-small" variant="text" class="ml-auto opacity-0 group-hover:opacity-100" @click="handleDelete(comment.id)" />
          </div>
          <p class="text-sm text-text-primary mt-0.5 leading-relaxed">{{ comment.content }}</p>
        </div>
      </div>
      <p v-if="comments.length === 0 && !loading" class="text-sm text-text-secondary text-center py-2">{{ t('task.noComments') }}</p>
    </div>
    <div class="flex gap-2">
      <v-text-field v-model="newComment" :placeholder="t('task.writeComment')" density="compact" variant="outlined" rounded="lg" hide-details @keyup.enter="handleSubmit" />
      <v-btn icon="mdi-send" color="primary" size="small" :loading="sending" :disabled="!newComment.trim()" @click="handleSubmit" />
    </div>
  </div>
</template>
