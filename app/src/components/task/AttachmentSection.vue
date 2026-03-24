<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { attachmentService, type Attachment } from '@/services/attachmentService'
import { useToast } from '@/composables/useToast'
import { formatRelativeTime } from '@/utils/date'

const { t } = useI18n()
const toast = useToast()
const props = defineProps<{ taskId: string }>()

const attachments = ref<Attachment[]>([])
const uploading = ref(false)

async function loadAttachments() {
  try { attachments.value = await attachmentService.list(props.taskId) }
  catch { /* */ }
}
onMounted(loadAttachments)

function triggerUpload() {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*,.pdf,.doc,.docx,.txt,.csv,.xlsx'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) handleFile(file)
  }
  input.click()
}

async function handleFile(file: File) {
  if (file.size > 5 * 1024 * 1024) {
    toast.error('Max 5MB')
    return
  }

  uploading.value = true
  try {
    const data = await readFileAsBase64(file)
    const attachment = await attachmentService.upload(props.taskId, {
      filename: file.name,
      contentType: file.type,
      size: file.size,
      data,
    })
    attachments.value.unshift(attachment)
    toast.success(file.name + ' ✓')
  } catch {
    toast.error('Upload failed')
  } finally {
    uploading.value = false
  }
}

function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

async function handleDownload(attachment: Attachment) {
  try {
    const full = await attachmentService.download(attachment.id)
    const link = document.createElement('a')
    link.href = full.data
    link.download = full.filename
    link.click()
  } catch {
    toast.error('Download failed')
  }
}

async function handleDelete(id: string) {
  await attachmentService.delete(id)
  attachments.value = attachments.value.filter((a) => a.id !== id)
  toast.success(t('common.delete') + ' ✓')
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + 'B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + 'KB'
  return (bytes / (1024 * 1024)).toFixed(1) + 'MB'
}

function isImage(contentType: string): boolean {
  return contentType.startsWith('image/')
}

const fileIcons: Record<string, string> = {
  'application/pdf': 'mdi-file-pdf-box',
  'text/plain': 'mdi-file-document-outline',
  'text/csv': 'mdi-file-delimited-outline',
}
</script>

<template>
  <div>
    <div class="flex items-center justify-between mb-2">
      <p class="text-[10px] font-semibold uppercase tracking-widest text-white/25">
        {{ t('task.attachments') ?? 'Attachments' }} ({{ attachments.length }})
      </p>
      <button
        class="flex items-center gap-1 px-2 py-1 rounded-lg text-[11px] text-primary-light hover:bg-primary/10 transition-colors"
        :disabled="uploading"
        @click="triggerUpload"
      >
        <v-icon :icon="uploading ? 'mdi-loading mdi-spin' : 'mdi-paperclip'" size="14" />
        {{ uploading ? '...' : t('common.add') }}
      </button>
    </div>

    <!-- Attachment list -->
    <div v-if="attachments.length > 0" class="space-y-1.5">
      <div
        v-for="a in attachments"
        :key="a.id"
        class="flex items-center gap-2.5 p-2 rounded-lg bg-white/[0.02] hover:bg-white/[0.04] transition-colors group"
      >
        <!-- Icon -->
        <div class="w-8 h-8 rounded-lg flex items-center justify-center flex-shrink-0" :class="isImage(a.contentType) ? 'bg-secondary/15' : 'bg-white/5'">
          <v-icon :icon="isImage(a.contentType) ? 'mdi-image-outline' : (fileIcons[a.contentType] ?? 'mdi-file-outline')" size="16" :class="isImage(a.contentType) ? 'text-secondary' : 'text-white/30'" />
        </div>

        <!-- Info -->
        <div class="flex-1 min-w-0 cursor-pointer" @click="handleDownload(a)">
          <p class="text-xs font-medium truncate hover:text-primary-light transition-colors">{{ a.filename }}</p>
          <p class="text-[10px] text-white/20">{{ formatSize(a.size) }} · {{ formatRelativeTime(a.createdAt) }}</p>
        </div>

        <!-- Actions -->
        <button class="p-1 rounded opacity-0 group-hover:opacity-100 hover:bg-white/5 transition-all" @click="handleDownload(a)">
          <v-icon icon="mdi-download-outline" size="14" class="text-white/30" />
        </button>
        <button class="p-1 rounded opacity-0 group-hover:opacity-100 hover:bg-error/10 transition-all" @click="handleDelete(a.id)">
          <v-icon icon="mdi-delete-outline" size="14" class="text-error/50" />
        </button>
      </div>
    </div>

    <p v-else class="text-[11px] text-white/15 text-center py-3">No files attached</p>
  </div>
</template>
