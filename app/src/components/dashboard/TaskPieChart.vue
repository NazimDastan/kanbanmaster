<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'
import { useTheme } from '@/composables/useTheme'

ChartJS.register(ArcElement, Tooltip, Legend)

const { theme } = useTheme()

const props = defineProps<{
  completed: number
  inProgress: number
  overdue: number
}>()

const chartData = computed(() => {
  const isDark = theme.value === 'dark'
  return {
    labels: ['Completed', 'In Progress', 'Overdue'],
    datasets: [{
      data: [props.completed, props.inProgress, props.overdue],
      backgroundColor: ['#10b981', '#6366f1', '#ef4444'],
      borderColor: isDark ? '#0a0a0f' : '#ffffff',
      borderWidth: 3,
      hoverOffset: 6,
    }],
  }
})

const chartOptions = computed(() => {
  const isDark = theme.value === 'dark'
  return {
    responsive: true,
    maintainAspectRatio: false,
    cutout: '70%',
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: isDark ? '#1e1e32' : '#ffffff',
        titleColor: isDark ? '#f1f5f9' : '#1e293b',
        bodyColor: isDark ? '#94a3b8' : '#64748b',
        borderColor: isDark ? 'rgba(255,255,255,0.05)' : 'rgba(0,0,0,0.1)',
        borderWidth: 1,
        cornerRadius: 8,
        padding: 10,
      },
    },
  }
})
</script>

<template>
  <div class="relative h-[180px]">
    <Doughnut :data="chartData" :options="chartOptions" />
    <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
      <div class="text-center">
        <p class="text-2xl font-bold text-[var(--text)]">{{ completed + inProgress + overdue }}</p>
        <p class="text-[10px] text-[var(--text-muted)]">Total</p>
      </div>
    </div>
  </div>
</template>
