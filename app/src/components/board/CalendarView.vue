<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Task, Priority } from '@/types/task'
import { PRIORITY_CONFIG } from '@/types/task'

useI18n()

const props = defineProps<{ tasks: Task[] }>()
const emit = defineEmits<{ taskClick: [task: Task] }>()

const today = new Date()
const currentMonth = ref(today.getMonth())
const currentYear = ref(today.getFullYear())

const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
const dayNames = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun']

const currentMonthName = computed(() => `${monthNames[currentMonth.value]} ${currentYear.value}`)

function prevMonth() {
  if (currentMonth.value === 0) { currentMonth.value = 11; currentYear.value-- }
  else currentMonth.value--
}

function nextMonth() {
  if (currentMonth.value === 11) { currentMonth.value = 0; currentYear.value++ }
  else currentMonth.value++
}

interface CalendarDay {
  date: number
  month: number
  year: number
  isCurrentMonth: boolean
  isToday: boolean
  tasks: Task[]
}

const calendarDays = computed<CalendarDay[]>(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)

  // Monday = 0 start
  let startDay = firstDay.getDay() - 1
  if (startDay < 0) startDay = 6

  const days: CalendarDay[] = []

  // Previous month fill
  const prevLastDay = new Date(year, month, 0)
  for (let i = startDay - 1; i >= 0; i--) {
    days.push({ date: prevLastDay.getDate() - i, month: month - 1, year, isCurrentMonth: false, isToday: false, tasks: [] })
  }

  // Current month
  for (let d = 1; d <= lastDay.getDate(); d++) {
    const isToday = d === today.getDate() && month === today.getMonth() && year === today.getFullYear()
    const dateStr = `${year}-${String(month + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
    const dayTasks = props.tasks.filter((task) => task.deadline?.startsWith(dateStr))
    days.push({ date: d, month, year, isCurrentMonth: true, isToday, tasks: dayTasks })
  }

  // Next month fill
  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    days.push({ date: i, month: month + 1, year, isCurrentMonth: false, isToday: false, tasks: [] })
  }

  return days
})

function priorityColor(p: Priority): string {
  return PRIORITY_CONFIG[p].color
}
</script>

<template>
  <div>
    <!-- Month nav -->
    <div class="flex items-center justify-between mb-4">
      <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" @click="prevMonth">
        <v-icon icon="mdi-chevron-left" size="18" class="text-white/40" />
      </button>
      <h3 class="text-sm font-semibold">{{ currentMonthName }}</h3>
      <button class="p-1.5 rounded-lg hover:bg-white/5 transition-colors" @click="nextMonth">
        <v-icon icon="mdi-chevron-right" size="18" class="text-white/40" />
      </button>
    </div>

    <!-- Day headers -->
    <div class="grid grid-cols-7 gap-px mb-1">
      <div v-for="day in dayNames" :key="day" class="text-center text-[10px] font-medium text-white/25 py-1">
        {{ day }}
      </div>
    </div>

    <!-- Calendar grid -->
    <div class="grid grid-cols-7 gap-px">
      <div
        v-for="(day, i) in calendarDays"
        :key="i"
        class="min-h-[70px] p-1 rounded-lg border transition-colors"
        :class="{
          'border-white/[0.04] bg-card': day.isCurrentMonth,
          'border-transparent bg-transparent': !day.isCurrentMonth,
          '!border-primary/30 !bg-primary/5': day.isToday,
        }"
      >
        <p class="text-[11px] mb-0.5" :class="day.isCurrentMonth ? (day.isToday ? 'text-primary-light font-bold' : 'text-white/50') : 'text-white/15'">
          {{ day.date }}
        </p>
        <div class="space-y-0.5">
          <button
            v-for="task in day.tasks.slice(0, 2)"
            :key="task.id"
            class="w-full text-left px-1 py-0.5 rounded text-[9px] truncate hover:brightness-125 transition-all"
            :style="{ backgroundColor: priorityColor(task.priority) + '20', color: priorityColor(task.priority) }"
            @click="emit('taskClick', task)"
          >
            {{ task.title }}
          </button>
          <p v-if="day.tasks.length > 2" class="text-[9px] text-white/25 px-1">+{{ day.tasks.length - 2 }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
