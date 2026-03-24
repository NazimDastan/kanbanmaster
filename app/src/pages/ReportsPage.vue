<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { reportService, type ReportRequest } from '@/services/reportService'
import StatCard from '@/components/common/StatCard.vue'
import { formatRelativeTime } from '@/utils/date'
import { getInitials } from '@/utils/format'
import PerformanceBarChart from '@/components/dashboard/PerformanceBarChart.vue'

const { t } = useI18n()

const activeTab = ref('performance')
const incomingReports = ref<ReportRequest[]>([])
const sentReports = ref<ReportRequest[]>([])
const loading = ref(false)
const showRequestModal = ref(false)
const showRespondModal = ref(false)
const respondingReport = ref<ReportRequest | null>(null)
const requestTargetUserId = ref('')
const requestTeamId = ref('')
const requestMessage = ref('')
const responseText = ref('')
const submitLoading = ref(false)

const performanceData = [
  { name: 'Alice Brown', score: 95, completed: 20, overdue: 0, onTime: 20 },
  { name: 'John Doe', score: 92, completed: 18, overdue: 1, onTime: 17 },
  { name: 'Jane Smith', score: 87, completed: 15, overdue: 2, onTime: 13 },
  { name: 'Bob Wilson', score: 78, completed: 12, overdue: 3, onTime: 9 },
]

const pendingCount = computed(() => incomingReports.value.filter(r => r.status === 'pending').length)
const statusColors: Record<string, string> = { pending: '#f59e0b', submitted: '#6366f1', reviewed: '#10b981' }

async function loadReports() {
  loading.value = true
  try {
    const [inc, sent] = await Promise.all([reportService.getIncoming(), reportService.getSent()])
    incomingReports.value = inc
    sentReports.value = sent
  } catch { /* API not connected */ }
  finally { loading.value = false }
}
onMounted(loadReports)

async function handleRequestReport() {
  if (!requestTargetUserId.value || !requestTeamId.value) return
  submitLoading.value = true
  try {
    await reportService.requestReport({ targetUserId: requestTargetUserId.value, teamId: requestTeamId.value, message: requestMessage.value })
    showRequestModal.value = false
    await loadReports()
  } finally { submitLoading.value = false }
}

function openRespond(report: ReportRequest) { respondingReport.value = report; responseText.value = ''; showRespondModal.value = true }

async function handleRespond() {
  if (!respondingReport.value || !responseText.value) return
  submitLoading.value = true
  try { await reportService.respond(respondingReport.value.id, responseText.value); showRespondModal.value = false; await loadReports() }
  finally { submitLoading.value = false }
}

async function handleReview(id: string) { await reportService.review(id); await loadReports() }

function scoreColor(s: number) { return s >= 90 ? '#10b981' : s >= 75 ? '#6366f1' : s >= 60 ? '#f59e0b' : '#ef4444' }
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 space-y-5">
    <div class="flex items-center justify-between">
      <h1 class="text-xl font-bold">{{ t('reports.title') }}</h1>
      <v-btn color="primary" prepend-icon="mdi-file-document-plus-outline" size="small" style="text-transform: none" @click="showRequestModal = true">{{ t('reports.requestReport') }}</v-btn>
    </div>

    <!-- Stats -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard :title="t('reports.pending')" :value="pendingCount" icon="mdi-clock-outline" gradient="linear-gradient(135deg, #f59e0b, #fbbf24)" />
      <StatCard :title="t('reports.efficiency')" value="85%" icon="mdi-speedometer" gradient="linear-gradient(135deg, #6366f1, #818cf8)" :change="5" />
      <StatCard :title="t('reports.completedLabel')" value="75" icon="mdi-check-all" gradient="linear-gradient(135deg, #10b981, #34d399)" :change="12" />
      <StatCard :title="t('reports.onTime')" value="89%" icon="mdi-clock-check-outline" gradient="linear-gradient(135deg, #a855f7, #c084fc)" />
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 p-0.5 bg-white/5 rounded-lg w-fit">
      <button v-for="tab in ['performance', 'incoming', 'sent']" :key="tab" class="px-3 py-1.5 rounded-md text-xs font-medium transition-all" :class="activeTab === tab ? 'bg-primary text-white' : 'text-white/40 hover:text-white/60'" @click="activeTab = tab">
        {{ tab === 'performance' ? t('reports.performance') : tab === 'incoming' ? t('reports.incoming') : t('reports.sent') }}
      </button>
    </div>

    <!-- Performance -->
    <div v-if="activeTab === 'performance'" class="rounded-xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
      <div class="grid grid-cols-[auto_1fr_auto_auto_auto_auto_1fr] gap-x-4 text-[11px] font-semibold uppercase tracking-widest text-white/30 px-4 py-2.5 border-b border-white/5">
        <span>#</span><span>{{ t('reports.member') }}</span><span>{{ t('reports.score') }}</span><span>{{ t('reports.done') }}</span><span>{{ t('reports.onTimeCol') }}</span><span>{{ t('reports.late') }}</span><span>{{ t('reports.progress') }}</span>
      </div>
      <div v-for="(m, i) in performanceData" :key="m.name" class="grid grid-cols-[auto_1fr_auto_auto_auto_auto_1fr] gap-x-4 items-center px-4 py-3 border-b border-white/[0.03] hover:bg-white/[0.02] transition-colors">
        <span class="text-xs font-medium text-white/30 w-4">{{ i + 1 }}</span>
        <div class="flex items-center gap-2.5">
          <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
            <span class="text-[9px] text-white font-semibold">{{ getInitials(m.name) }}</span>
          </div>
          <span class="text-sm font-medium truncate">{{ m.name }}</span>
        </div>
        <span class="px-2 py-0.5 rounded-full text-[11px] font-semibold" :style="{ backgroundColor: scoreColor(m.score) + '18', color: scoreColor(m.score) }">{{ m.score }}%</span>
        <span class="text-xs text-white/60 w-8 text-center">{{ m.completed }}</span>
        <span class="text-xs text-success w-10 text-center">{{ m.onTime }}</span>
        <span class="text-xs text-error w-8 text-center">{{ m.overdue }}</span>
        <v-progress-linear :model-value="m.score" :color="scoreColor(m.score)" height="5" rounded bg-color="#1e1e32" />
      </div>
    </div>

    <!-- Performance Chart -->
    <div v-if="activeTab === 'performance'" class="rounded-2xl border border-white/5 bg-[#0f0f1a] p-5 mt-4">
      <h3 class="text-sm font-semibold mb-4">{{ t('reports.onTimeCol') }} vs {{ t('reports.late') }}</h3>
      <PerformanceBarChart :members="performanceData" />
      <div class="flex justify-center gap-4 mt-4">
        <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-success" /><span class="text-[11px] text-white/40">{{ t('reports.onTimeCol') }}</span></div>
        <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-error" /><span class="text-[11px] text-white/40">{{ t('reports.late') }}</span></div>
      </div>
    </div>

    <!-- Incoming / Sent reports (same structure) -->
    <div v-if="activeTab === 'incoming' || activeTab === 'sent'" class="rounded-xl border border-white/5 bg-[#0f0f1a] overflow-hidden">
      <div v-for="report in (activeTab === 'incoming' ? incomingReports : sentReports)" :key="report.id" class="flex items-center gap-3 px-4 py-3 border-b border-white/[0.03]">
        <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary/50 to-secondary/50 flex items-center justify-center flex-shrink-0">
          <span class="text-[9px] text-white font-semibold">{{ getInitials((activeTab === 'incoming' ? report.requester?.name : report.targetUser?.name) ?? 'U') }}</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm truncate">{{ activeTab === 'incoming' ? t('reports.requestedReport', { name: report.requester?.name ?? t('reports.unknown') }) : t('reports.reportFrom', { name: report.targetUser?.name ?? t('reports.unknown') }) }}</p>
          <p class="text-[11px] text-white/25">{{ report.message || t('reports.noMessage') }} · {{ formatRelativeTime(report.createdAt) }}</p>
        </div>
        <span class="px-2 py-0.5 rounded-full text-[10px] font-medium" :style="{ backgroundColor: (statusColors[report.status] ?? '#64748b') + '18', color: statusColors[report.status] }">{{ report.status }}</span>
        <v-btn v-if="activeTab === 'incoming' && report.status === 'pending'" size="x-small" variant="tonal" color="primary" style="text-transform: none" @click="openRespond(report)">{{ t('reports.respond') }}</v-btn>
        <v-btn v-if="activeTab === 'sent' && report.status === 'submitted'" size="x-small" variant="tonal" color="success" style="text-transform: none" @click="handleReview(report.id)">{{ t('reports.review') }}</v-btn>
      </div>
      <div v-if="(activeTab === 'incoming' ? incomingReports : sentReports).length === 0" class="py-8 text-center text-xs text-white/25">{{ t('reports.noReports') }}</div>
    </div>

    <!-- Request modal -->
    <v-dialog v-model="showRequestModal" max-width="420">
      <v-card class="pa-5" color="#161625">
        <h3 class="text-sm font-bold mb-3">{{ t('reports.requestReport') }}</h3>
        <v-text-field v-model="requestTargetUserId" :label="t('reports.targetUserId')" prepend-inner-icon="mdi-account-outline" class="mb-2" />
        <v-text-field v-model="requestTeamId" :label="t('reports.teamId')" prepend-inner-icon="mdi-account-group-outline" class="mb-2" />
        <v-textarea v-model="requestMessage" :label="t('reports.message')" rows="3" auto-grow variant="outlined" rounded="lg" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showRequestModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :loading="submitLoading" :disabled="!requestTargetUserId || !requestTeamId" style="text-transform: none" @click="handleRequestReport">{{ t('reports.sendRequest') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <v-dialog v-model="showRespondModal" max-width="420">
      <v-card class="pa-5" color="#161625">
        <h3 class="text-sm font-bold mb-2">{{ t('reports.respond') }}</h3>
        <p class="text-xs text-white/30 mb-3">{{ respondingReport?.message }}</p>
        <v-textarea v-model="responseText" :label="t('reports.yourResponse')" rows="4" auto-grow variant="outlined" rounded="lg" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showRespondModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :loading="submitLoading" :disabled="!responseText" style="text-transform: none" @click="handleRespond">{{ t('common.submit') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
