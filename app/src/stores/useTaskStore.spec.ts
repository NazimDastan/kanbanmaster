// @vitest-environment jsdom
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTaskStore } from './useTaskStore'

vi.mock('@/services/taskService', () => {
  const base = {
    id: 't1', columnId: 'c1', title: 'Test Task', description: 'desc',
    creatorId: 'u1', assigneeId: null, assignee: null, priority: 'medium',
    deadline: null, position: 0, subtasks: [], labels: [],
    createdAt: '', updatedAt: '', completedAt: null,
  }
  return {
    taskService: {
      get: vi.fn().mockResolvedValue({ ...base }),
      create: vi.fn().mockResolvedValue({ ...base, id: 't2', title: 'New Task' }),
      update: vi.fn().mockResolvedValue({ ...base, title: 'Updated' }),
      delete: vi.fn().mockResolvedValue(undefined),
      move: vi.fn().mockResolvedValue(undefined),
      assign: vi.fn().mockResolvedValue(undefined),
    },
  }
})

describe('useTaskStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with no selected task', () => {
    const store = useTaskStore()
    expect(store.selectedTask).toBeNull()
    expect(store.loading).toBe(false)
  })

  it('fetchTask loads task', async () => {
    const store = useTaskStore()
    await store.fetchTask('t1')
    expect(store.selectedTask).not.toBeNull()
    expect(store.selectedTask?.title).toBe('Test Task')
  })

  it('createTask returns new task', async () => {
    const store = useTaskStore()
    const task = await store.createTask({
      columnId: 'c1', title: 'New Task', priority: 'high',
    })
    expect(task.title).toBe('New Task')
  })

  it('updateTask updates selected task', async () => {
    const store = useTaskStore()
    await store.fetchTask('t1')
    const updated = await store.updateTask('t1', { title: 'Updated' })
    expect(updated.title).toBe('Updated')
    expect(store.selectedTask?.title).toBe('Updated')
  })

  it('deleteTask clears selected task', async () => {
    const store = useTaskStore()
    await store.fetchTask('t1')
    await store.deleteTask('t1')
    expect(store.selectedTask).toBeNull()
  })

  it('moveTask calls service', async () => {
    const store = useTaskStore()
    await expect(store.moveTask('t1', 'c2', 0)).resolves.toBeUndefined()
  })

  it('assignTask calls service', async () => {
    const store = useTaskStore()
    await expect(store.assignTask('t1', 'u2')).resolves.toBeUndefined()
  })
})
