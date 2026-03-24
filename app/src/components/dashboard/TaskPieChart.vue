<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js'

ChartJS.register(ArcElement, Tooltip, Legend)

const props = defineProps<{
  completed: number
  inProgress: number
  overdue: number
}>()

const chartData = computed(() => ({
  labels: ['Completed', 'In Progress', 'Overdue'],
  datasets: [{
    data: [props.completed, props.inProgress, props.overdue],
    backgroundColor: ['#10b981', '#6366f1', '#ef4444'],
    borderColor: ['#0a0a0f', '#0a0a0f', '#0a0a0f'],
    borderWidth: 3,
    hoverOffset: 6,
  }],
}))

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  cutout: '70%',
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: '#1e1e32',
      titleColor: '#f1f5f9',
      bodyColor: '#94a3b8',
      borderColor: 'rgba(255,255,255,0.05)',
      borderWidth: 1,
      cornerRadius: 8,
      padding: 10,
    },
  },
}
</script>

<template>
  <div class="relative" style="height: 180px">
    <Doughnut :data="chartData" :options="chartOptions" />
    <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
      <div class="text-center">
        <p class="text-2xl font-bold">{{ completed + inProgress + overdue }}</p>
        <p class="text-[10px] text-white/30">Total</p>
      </div>
    </div>
  </div>
</template>
