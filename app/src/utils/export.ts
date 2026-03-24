import type { Task } from '@/types/task'

export function exportTasksToCSV(tasks: Task[], filename = 'tasks.csv') {
  const headers = ['Title', 'Priority', 'Assignee', 'Deadline', 'Status', 'Created']
  const rows = tasks.map((t) => [
    `"${t.title.replace(/"/g, '""')}"`,
    t.priority,
    t.assignee?.name ?? 'Unassigned',
    t.deadline?.slice(0, 10) ?? '',
    t.completedAt ? 'Completed' : 'Active',
    t.createdAt.slice(0, 10),
  ])

  const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n')
  const blob = new Blob(['\ufeff' + csv], { type: 'text/csv;charset=utf-8;' })
  const url = URL.createObjectURL(blob)

  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}
