<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { reportService, type ReportRequest } from '@/services/reportService'
import { dashboardService, type UserPerformance } from '@/services/dashboardService'
import { useTeamStore } from '@/stores/useTeamStore'
import { useToast } from '@/composables/useToast'
import StatCard from '@/components/common/StatCard.vue'
import { formatRelativeTime } from '@/utils/date'
import { getInitials } from '@/utils/format'
import { exportToCSV } from '@/utils/export'
import PerformanceBarChart from '@/components/dashboard/PerformanceBarChart.vue'

const { t } = useI18n()
const toast = useToast()
const teamStore = useTeamStore()

const activeTab = ref('performance')
const selectedTeamId = ref('')
const performanceData = ref<UserPerformance[]>([])
const performanceLoading = ref(false)

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

const pendingCount = computed(() => incomingReports.value.filter(r => r.status === 'pending').length)
const totalCompleted = computed(() => performanceData.value.reduce((s, m) => s + m.completed, 0))
const totalOnTime = computed(() => performanceData.value.reduce((s, m) => s + m.onTime, 0))
const avgScore = computed(() => {
  if (performanceData.value.length === 0) return 0
  return Math.round(performanceData.value.reduce((s, m) => s + m.score, 0) / performanceData.value.length)
})
const statusColors: Record<string, string> = { pending: '#f59e0b', submitted: '#6366f1', reviewed: '#10b981' }

async function loadData() {
  loading.value = true
  try {
    await teamStore.fetchTeams()
    if (teamStore.teams.length > 0 && !selectedTeamId.value) {
      selectedTeamId.value = teamStore.teams[0].id
    }
    const [inc, sent] = await Promise.all([reportService.getIncoming(), reportService.getSent()])
    incomingReports.value = inc
    sentReports.value = sent
  } catch { toast.error(t('common.error')) }
  finally { loading.value = false }
}

async function loadPerformance() {
  if (!selectedTeamId.value) return
  performanceLoading.value = true
  try {
    performanceData.value = await dashboardService.getTeamPerformance(selectedTeamId.value)
  } catch {
    performanceData.value = []
    toast.error(t('common.error'))
  } finally {
    performanceLoading.value = false
  }
}

watch(selectedTeamId, loadPerformance)

async function init() {
  await loadData()
  await loadPerformance()
}
onMounted(init)

async function handleRequestReport() {
  if (!requestTargetUserId.value || !requestTeamId.value) return
  submitLoading.value = true
  try {
    await reportService.requestReport({ targetUserId: requestTargetUserId.value, teamId: requestTeamId.value, message: requestMessage.value })
    showRequestModal.value = false
    await loadData()
  } finally { submitLoading.value = false }
}

function openRespond(report: ReportRequest) { respondingReport.value = report; responseText.value = ''; showRespondModal.value = true }

async function handleRespond() {
  if (!respondingReport.value || !responseText.value) return
  submitLoading.value = true
  try { await reportService.respond(respondingReport.value.id, responseText.value); showRespondModal.value = false; await loadData() }
  finally { submitLoading.value = false }
}

async function handleReview(id: string) { await reportService.review(id); await loadData() }

function scoreColor(s: number) { return s >= 90 ? '#10b981' : s >= 75 ? '#6366f1' : s >= 60 ? '#f59e0b' : '#ef4444' }

function handleExportCSV() {
  if (activeTab.value === 'performance' && performanceData.value.length > 0) {
    const rows = performanceData.value.map((m) => ({
      Member: m.userName,
      Score: m.score,
      Completed: m.completed,
      'On Time': m.onTime,
      Late: m.overdue,
      'Total Tasks': m.totalTasks,
    }))
    exportToCSV(rows, 'performance-report.csv')
  } else if (activeTab.value === 'incoming' && incomingReports.value.length > 0) {
    const rows = incomingReports.value.map((r) => ({
      Requester: r.requester?.name ?? '',
      Status: r.status,
      Message: r.message,
      'Created At': r.createdAt,
      Response: r.response ?? '',
    }))
    exportToCSV(rows, 'incoming-reports.csv')
  } else if (activeTab.value === 'sent' && sentReports.value.length > 0) {
    const rows = sentReports.value.map((r) => ({
      'Target User': r.targetUser?.name ?? '',
      Status: r.status,
      Message: r.message,
      'Created At': r.createdAt,
      Response: r.response ?? '',
    }))
    exportToCSV(rows, 'sent-reports.csv')
  }
}

const canExport = computed(() => {
  if (activeTab.value === 'performance') return performanceData.value.length > 0
  if (activeTab.value === 'incoming') return incomingReports.value.length > 0
  if (activeTab.value === 'sent') return sentReports.value.length > 0
  return false
})
</script>

<template>
  <div class="p-4 md:p-6 lg:p-8 space-y-5">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
      <h1 class="text-xl font-bold">{{ t('reports.title') }}</h1>
      <div class="flex items-center gap-2">
        <!-- Team selector -->
        <select
          v-model="selectedTeamId"
          class="px-3 py-2 rounded-lg text-sm outline-none bg-[var(--bg-input)] border border-[var(--border)] text-[var(--text)]"
        >
          <option v-for="team in teamStore.teams" :key="team.id" :value="team.id">{{ team.name }}</option>
        </select>
        <v-btn color="primary" prepend-icon="mdi-file-document-plus-outline" size="small" @click="showRequestModal = true">{{ t('reports.requestReport') }}</v-btn>
        <v-btn variant="outlined" prepend-icon="mdi-download" size="small" :disabled="!canExport" @click="handleExportCSV">{{ t('reports.exportCSV') }}</v-btn>
      </div>
    </div>

    <!-- Stats (from real data) -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <StatCard :title="t('reports.pending')" :value="pendingCount" icon="mdi-clock-outline" gradient="linear-gradient(135deg, #f59e0b, #fbbf24)" />
      <StatCard :title="t('reports.efficiency')" :value="avgScore + '%'" icon="mdi-speedometer" gradient="linear-gradient(135deg, #6366f1, #818cf8)" />
      <StatCard :title="t('reports.completedLabel')" :value="totalCompleted" icon="mdi-check-all" gradient="linear-gradient(135deg, #10b981, #34d399)" />
      <StatCard :title="t('reports.onTime')" :value="totalOnTime" icon="mdi-clock-check-outline" gradient="linear-gradient(135deg, #a855f7, #c084fc)" />
    </div>

    <!-- Tabs -->
    <div class="flex gap-1 p-0.5 rounded-lg w-fit bg-[var(--bg-input)]">
      <button v-for="tab in (['performance', 'incoming', 'sent'] as const)" :key="tab" class="px-3 py-1.5 rounded-md text-xs font-medium transition-all" :class="activeTab === tab ? 'bg-primary text-white' : ''" :style="activeTab !== tab ? { color: 'var(--text-secondary)' } : {}" @click="activeTab = tab">
        {{ tab === 'performance' ? t('reports.performance') : tab === 'incoming' ? t('reports.incoming') : t('reports.sent') }}
      </button>
    </div>

    <!-- Performance -->
    <template v-if="activeTab === 'performance'">
      <div v-if="performanceLoading" class="flex justify-center py-12">
        <v-progress-circular indeterminate color="primary" size="40" />
      </div>

      <div v-else-if="performanceData.length === 0" class="rounded-xl p-8 text-center bg-[var(--bg-card)] border border-[var(--border)]">
        <v-icon icon="mdi-chart-bar" size="40" class="text-[var(--text-faint)]" />
        <p class="text-sm mt-3 text-[var(--text-muted)]">{{ t('common.noData') }}</p>
        <p class="text-xs mt-1 text-[var(--text-faint)]">{{ teamStore.teams.length === 0 ? t('team.noTeamsYet') : t('dashboard.createTasksHint') }}</p>
      </div>

      <template v-else>
        <!-- Table -->
        <div class="rounded-xl overflow-hidden overflow-x-auto bg-[var(--bg-card)] border border-[var(--border)]">
          <div class="min-w-[600px]">
          <div class="grid grid-cols-[auto_1fr_auto_auto_auto_auto_1fr] gap-x-4 text-[10px] font-semibold uppercase tracking-widest px-4 py-2.5 text-[var(--text-muted)] border-b border-[var(--border)]">
            <span>#</span><span>{{ t('reports.member') }}</span><span>{{ t('reports.score') }}</span><span>{{ t('reports.done') }}</span><span>{{ t('reports.onTimeCol') }}</span><span>{{ t('reports.late') }}</span><span>{{ t('reports.progress') }}</span>
          </div>
          <div v-for="(m, i) in performanceData" :key="m.userId" class="grid grid-cols-[auto_1fr_auto_auto_auto_auto_1fr] gap-x-4 items-center px-4 py-3 border-b border-[var(--border)]">
            <span class="text-xs font-medium w-4 text-[var(--text-muted)]">{{ i + 1 }}</span>
            <div class="flex items-center gap-2.5">
              <div class="w-7 h-7 rounded-full bg-gradient-to-br from-primary/60 to-secondary/60 flex items-center justify-center flex-shrink-0">
                <span class="text-[9px] text-white font-semibold">{{ getInitials(m.userName) }}</span>
              </div>
              <span class="text-sm font-medium truncate text-[var(--text)]">{{ m.userName }}</span>
            </div>
            <span class="px-2 py-0.5 rounded-full text-[11px] font-semibold" :style="{ backgroundColor: scoreColor(m.score) + '18', color: scoreColor(m.score) }">{{ m.score }}%</span>
            <span class="text-xs w-8 text-center text-[var(--text-secondary)]">{{ m.completed }}</span>
            <span class="text-xs text-success w-10 text-center">{{ m.onTime }}</span>
            <span class="text-xs text-error w-8 text-center">{{ m.overdue }}</span>
            <v-progress-linear :model-value="m.score" :color="scoreColor(m.score)" height="5" rounded bg-color="surface" />
          </div>
          </div>
        </div>

        <!-- Chart -->
        <div v-if="performanceData.length > 0" class="rounded-2xl p-5 bg-[var(--bg-card)] border border-[var(--border)]">
          <h3 class="text-sm font-semibold mb-4 text-[var(--text)]">{{ t('reports.onTimeCol') }} vs {{ t('reports.late') }}</h3>
          <PerformanceBarChart :members="performanceData" />
          <div class="flex justify-center gap-4 mt-4">
            <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-success" /><span class="text-[11px] text-[var(--text-secondary)]">{{ t('reports.onTimeCol') }}</span></div>
            <div class="flex items-center gap-1.5"><div class="w-2.5 h-2.5 rounded-full bg-error" /><span class="text-[11px] text-[var(--text-secondary)]">{{ t('reports.late') }}</span></div>
          </div>
        </div>
      </template>
    </template>

    <!-- Incoming / Sent -->
    <div v-if="activeTab === 'incoming' || activeTab === 'sent'" class="rounded-xl overflow-hidden bg-[var(--bg-card)] border border-[var(--border)]">
      <div v-for="report in (activeTab === 'incoming' ? incomingReports : sentReports)" :key="report.id" class="flex items-center gap-3 px-4 py-3 border-b border-[var(--border)]">
        <div class="w-8 h-8 rounded-full bg-gradient-to-br from-primary/50 to-secondary/50 flex items-center justify-center flex-shrink-0">
          <span class="text-[9px] text-white font-semibold">{{ getInitials((activeTab === 'incoming' ? report.requester?.name : report.targetUser?.name) ?? 'U') }}</span>
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm truncate text-[var(--text)]">{{ activeTab === 'incoming' ? t('reports.requestedReport', { name: report.requester?.name ?? t('reports.unknown') }) : t('reports.reportFrom', { name: report.targetUser?.name ?? t('reports.unknown') }) }}</p>
          <p class="text-[11px] text-[var(--text-muted)]">{{ report.message || t('reports.noMessage') }} · {{ formatRelativeTime(report.createdAt) }}</p>
        </div>
        <span class="px-2 py-0.5 rounded-full text-[10px] font-medium" :style="{ backgroundColor: (statusColors[report.status] ?? '#64748b') + '18', color: statusColors[report.status] }">{{ report.status }}</span>
        <v-btn v-if="activeTab === 'incoming' && report.status === 'pending'" size="x-small" variant="tonal" color="primary" @click="openRespond(report)">{{ t('reports.respond') }}</v-btn>
        <v-btn v-if="activeTab === 'sent' && report.status === 'submitted'" size="x-small" variant="tonal" color="success" @click="handleReview(report.id)">{{ t('reports.review') }}</v-btn>
      </div>
      <div v-if="(activeTab === 'incoming' ? incomingReports : sentReports).length === 0" class="py-8 text-center text-xs text-[var(--text-muted)]">{{ t('reports.noReports') }}</div>
    </div>

    <!-- Request modal -->
    <v-dialog v-model="showRequestModal" max-width="420">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-3">{{ t('reports.requestReport') }}</h3>
        <v-text-field v-model="requestTargetUserId" :label="t('reports.targetUserId')" prepend-inner-icon="mdi-account-outline" class="mb-2" />
        <v-text-field v-model="requestTeamId" :label="t('reports.teamId')" prepend-inner-icon="mdi-account-group-outline" class="mb-2" />
        <v-textarea v-model="requestMessage" :label="t('reports.message')" rows="3" auto-grow variant="outlined" rounded="lg" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showRequestModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :loading="submitLoading" :disabled="!requestTargetUserId || !requestTeamId" @click="handleRequestReport">{{ t('reports.sendRequest') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>

    <!-- Respond modal -->
    <v-dialog v-model="showRespondModal" max-width="420">
      <v-card class="pa-5" color="surface">
        <h3 class="text-sm font-bold mb-2">{{ t('reports.respond') }}</h3>
        <p class="text-xs mb-3 text-[var(--text-muted)]">{{ respondingReport?.message }}</p>
        <v-textarea v-model="responseText" :label="t('reports.yourResponse')" rows="4" auto-grow variant="outlined" rounded="lg" class="mb-3" />
        <div class="flex justify-end gap-2">
          <v-btn variant="text" size="small" @click="showRespondModal = false">{{ t('common.cancel') }}</v-btn>
          <v-btn color="primary" size="small" :loading="submitLoading" :disabled="!responseText" @click="handleRespond">{{ t('common.submit') }}</v-btn>
        </div>
      </v-card>
    </v-dialog>
  </div>
</template>
