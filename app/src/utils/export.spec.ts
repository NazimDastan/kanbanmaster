import { describe, it, expect, vi, beforeEach } from 'vitest'
import { exportToCSV, exportTasksToCSV } from './export'
import type { Task } from '@/types/task'

let capturedContent = ''

class MockBlob {
  constructor(parts: BlobPart[]) {
    capturedContent = parts[0] as string
  }
}

describe('exportToCSV', () => {
  let clickSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    capturedContent = ''
    clickSpy = vi.fn()
    vi.stubGlobal('Blob', MockBlob)
    vi.spyOn(document, 'createElement').mockReturnValue({
      href: '',
      download: '',
      click: clickSpy,
    } as unknown as HTMLAnchorElement)
    vi.spyOn(URL, 'createObjectURL').mockReturnValue('blob:test')
    vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {})
  })

  it('generates correct CSV content from data', () => {
    // Arrange
    const data = [
      { name: 'Alice', age: '30' },
      { name: 'Bob', age: '25' },
    ]

    // Act
    exportToCSV(data, 'test.csv')

    // Assert
    expect(capturedContent).toContain('name,age')
    expect(capturedContent).toContain('"Alice","30"')
    expect(capturedContent).toContain('"Bob","25"')
    expect(clickSpy).toHaveBeenCalled()
  })

  it('does nothing with empty array', () => {
    // Arrange & Act
    exportToCSV([], 'empty.csv')

    // Assert
    expect(clickSpy).not.toHaveBeenCalled()
  })

  it('escapes double quotes in values', () => {
    // Arrange
    const data = [{ text: 'He said "hello"' }]

    // Act
    exportToCSV(data, 'quotes.csv')

    // Assert
    expect(capturedContent).toContain('"He said ""hello"""')
  })

  it('handles null and undefined values', () => {
    // Arrange
    const data = [{ a: null, b: undefined }] as Record<string, unknown>[]

    // Act
    exportToCSV(data, 'nulls.csv')

    // Assert
    expect(capturedContent).toContain('"",""')
  })
})

describe('exportTasksToCSV', () => {
  let clickSpy: ReturnType<typeof vi.fn>

  beforeEach(() => {
    capturedContent = ''
    clickSpy = vi.fn()
    vi.stubGlobal('Blob', MockBlob)
    vi.spyOn(document, 'createElement').mockReturnValue({
      href: '',
      download: '',
      click: clickSpy,
    } as unknown as HTMLAnchorElement)
    vi.spyOn(URL, 'createObjectURL').mockReturnValue('blob:test')
    vi.spyOn(URL, 'revokeObjectURL').mockImplementation(() => {})
  })

  it('exports tasks with correct headers and values', () => {
    // Arrange
    const tasks: Task[] = [
      {
        id: 't1',
        columnId: 'c1',
        title: 'Fix bug',
        description: null,
        creatorId: 'u1',
        assigneeId: 'u2',
        assignee: { id: 'u2', name: 'Jane', email: 'j@e.com', avatarUrl: null, createdAt: '', updatedAt: '' },
        priority: 'high',
        deadline: '2026-06-15T00:00:00Z',
        position: 0,
        subtasks: [],
        labels: [],
        createdAt: '2026-01-10T00:00:00Z',
        updatedAt: '2026-01-10T00:00:00Z',
        completedAt: null,
      },
    ]

    // Act
    exportTasksToCSV(tasks)

    // Assert
    expect(capturedContent).toContain('Title,Priority,Assignee,Deadline,Status,Created')
    expect(capturedContent).toContain('"Fix bug"')
    expect(capturedContent).toContain('high')
    expect(capturedContent).toContain('Jane')
    expect(capturedContent).toContain('Active')
  })
})
