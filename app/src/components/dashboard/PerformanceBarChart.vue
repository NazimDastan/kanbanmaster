<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Tooltip } from 'chart.js'
import { useTheme } from '@/composables/useTheme'

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip)

const { theme } = useTheme()

const props = defineProps<{
  members: { name?: string; userName?: string; onTime: number; overdue: number }[]
}>()

const chartData = computed(() => ({
  labels: props.members.map(m => (m.name ?? m.userName ?? '?').split(' ')[0]),
  datasets: [
    { label: 'On Time', data: props.members.map(m => m.onTime), backgroundColor: '#10b981', borderRadius: 4, barPercentage: 0.6 },
    { label: 'Overdue', data: props.members.map(m => m.overdue), backgroundColor: '#ef4444', borderRadius: 4, barPercentage: 0.6 },
  ],
}))

const chartOptions = computed(() => {
  const isDark = theme.value === 'dark'
  return {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        grid: { display: false },
        ticks: { color: isDark ? 'rgba(255,255,255,0.3)' : 'rgba(0,0,0,0.4)', font: { size: 11 } },
      },
      y: {
        grid: { color: isDark ? 'rgba(255,255,255,0.03)' : 'rgba(0,0,0,0.06)' },
        ticks: { color: isDark ? 'rgba(255,255,255,0.2)' : 'rgba(0,0,0,0.3)', font: { size: 11 } },
        beginAtZero: true,
      },
    },
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: isDark ? '#1e1e32' : '#ffffff',
        titleColor: isDark ? '#f1f5f9' : '#1e293b',
        bodyColor: isDark ? '#94a3b8' : '#64748b',
        borderColor: isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.1)',
        borderWidth: 1,
        cornerRadius: 8,
      },
    },
  }
})
</script>

<template>
  <div class="h-[220px]">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
