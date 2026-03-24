import { describe, it, expect } from 'vitest'
import { formatRelativeTime, formatDate, formatDeadline } from './date'

describe('formatRelativeTime', () => {
  it('returns "Just now" for recent dates', () => {
    const now = new Date().toISOString()
    expect(formatRelativeTime(now)).toBe('Just now')
  })

  it('returns minutes ago', () => {
    const date = new Date(Date.now() - 5 * 60 * 1000).toISOString()
    expect(formatRelativeTime(date)).toBe('5m ago')
  })

  it('returns hours ago', () => {
    const date = new Date(Date.now() - 3 * 60 * 60 * 1000).toISOString()
    expect(formatRelativeTime(date)).toBe('3h ago')
  })

  it('returns days ago', () => {
    const date = new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    expect(formatRelativeTime(date)).toBe('2d ago')
  })
})

describe('formatDate', () => {
  it('formats a date string', () => {
    const result = formatDate('2026-03-15T10:00:00Z')
    expect(result).toContain('Mar')
    expect(result).toContain('15')
  })
})

describe('formatDeadline', () => {
  it('marks past deadlines as overdue', () => {
    const yesterday = new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
    const result = formatDeadline(yesterday)
    expect(result.isOverdue).toBe(true)
    expect(result.text).toContain('overdue')
  })

  it('marks today deadlines as urgent', () => {
    const today = new Date(Date.now() + 1 * 60 * 60 * 1000).toISOString()
    const result = formatDeadline(today)
    expect(result.isUrgent).toBe(true)
  })

  it('returns normal format for far deadlines', () => {
    const future = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString()
    const result = formatDeadline(future)
    expect(result.isOverdue).toBe(false)
    expect(result.isUrgent).toBe(false)
  })
})
