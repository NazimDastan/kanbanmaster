// @vitest-environment jsdom
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useBoardStore } from './useBoardStore'

vi.mock('@/services/boardService', () => ({
  boardService: {
    list: vi.fn().mockResolvedValue([
      { id: 'b1', name: 'Sprint Board', teamId: 't1', createdAt: '' },
      { id: 'b2', name: 'Backlog', teamId: 't1', createdAt: '' },
    ]),
    get: vi.fn().mockResolvedValue({
      id: 'b1', name: 'Sprint Board', teamId: 't1', createdAt: '',
      columns: [
        { id: 'c1', boardId: 'b1', name: 'Todo', position: 0, tasks: [], createdAt: '' },
        { id: 'c2', boardId: 'b1', name: 'In Progress', position: 1, tasks: [], createdAt: '' },
        { id: 'c3', boardId: 'b1', name: 'Done', position: 2, tasks: [], createdAt: '' },
      ],
    }),
    create: vi.fn().mockResolvedValue({ id: 'b3', name: 'New Board', teamId: 't1', createdAt: '' }),
    delete: vi.fn().mockResolvedValue(undefined),
  },
}))

describe('useBoardStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts empty', () => {
    const store = useBoardStore()
    expect(store.boards).toEqual([])
    expect(store.currentBoard).toBeNull()
  })

  it('fetchBoards loads boards', async () => {
    const store = useBoardStore()
    await store.fetchBoards()
    expect(store.boards).toHaveLength(2)
  })

  it('fetchBoard loads board with columns', async () => {
    const store = useBoardStore()
    await store.fetchBoard('b1')
    expect(store.currentBoard).not.toBeNull()
    expect(store.currentBoard?.columns).toHaveLength(3)
    expect(store.currentBoard?.columns[0].name).toBe('Todo')
  })

  it('createBoard adds to list', async () => {
    const store = useBoardStore()
    const board = await store.createBoard('New Board', 't1')
    expect(board.name).toBe('New Board')
    expect(store.boards).toHaveLength(1)
  })

  it('deleteBoard removes from list and clears current', async () => {
    const store = useBoardStore()
    await store.fetchBoards()
    await store.fetchBoard('b1')
    await store.deleteBoard('b1')
    expect(store.boards).toHaveLength(1)
    expect(store.currentBoard).toBeNull()
  })
})
