export function formatRelativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now.getTime() - date.getTime()
  const diffSec = Math.floor(diffMs / 1000)
  const diffMin = Math.floor(diffSec / 60)
  const diffHour = Math.floor(diffMin / 60)
  const diffDay = Math.floor(diffHour / 24)

  if (diffSec < 60) return 'Just now'
  if (diffMin < 60) return `${diffMin}m ago`
  if (diffHour < 24) return `${diffHour}h ago`
  if (diffDay < 7) return `${diffDay}d ago`
  return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
}

export function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

export function formatDeadline(dateStr: string): { text: string; isOverdue: boolean; isUrgent: boolean } {
  const deadline = new Date(dateStr)
  const now = new Date()
  const diffMs = deadline.getTime() - now.getTime()
  const diffDays = Math.ceil(diffMs / (1000 * 60 * 60 * 24))

  if (diffDays < 0) return { text: `${Math.abs(diffDays)}d overdue`, isOverdue: true, isUrgent: false }
  if (diffDays === 0) return { text: 'Due today', isOverdue: false, isUrgent: true }
  if (diffDays === 1) return { text: 'Due tomorrow', isOverdue: false, isUrgent: true }
  if (diffDays <= 3) return { text: `${diffDays} days left`, isOverdue: false, isUrgent: true }
  return { text: formatDate(dateStr), isOverdue: false, isUrgent: false }
}
