<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Tooltip } from 'chart.js'

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip)

const props = defineProps<{
  members: { name: string; onTime: number; overdue: number }[]
}>()

const chartData = computed(() => ({
  labels: props.members.map(m => m.name.split(' ')[0]),
  datasets: [
    {
      label: 'On Time',
      data: props.members.map(m => m.onTime),
      backgroundColor: '#10b981',
      borderRadius: 4,
      barPercentage: 0.6,
    },
    {
      label: 'Overdue',
      data: props.members.map(m => m.overdue),
      backgroundColor: '#ef4444',
      borderRadius: 4,
      barPercentage: 0.6,
    },
  ],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    x: {
      grid: { display: false },
      ticks: { color: 'rgba(255,255,255,0.3)', font: { size: 11 } },
    },
    y: {
      grid: { color: 'rgba(255,255,255,0.03)' },
      ticks: { color: 'rgba(255,255,255,0.2)', font: { size: 11 } },
      beginAtZero: true,
    },
  },
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#1e1e32',
      titleColor: '#f1f5f9',
      bodyColor: '#94a3b8',
      borderColor: 'rgba(255,255,255,0.05)',
      borderWidth: 1,
      cornerRadius: 8,
    },
  },
}
</script>

<template>
  <div style="height: 220px">
    <Bar :data="chartData" :options="chartOptions" />
  </div>
</template>
