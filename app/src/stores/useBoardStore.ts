import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Board, BoardWithColumns } from '@/types/board'
import { boardService } from '@/services/boardService'

export const useBoardStore = defineStore('board', () => {
  const boards = ref<Board[]>([])
  const currentBoard = ref<BoardWithColumns | null>(null)
  const loading = ref(false)

  async function fetchBoards() {
    loading.value = true
    try {
      boards.value = await boardService.list()
    } finally {
      loading.value = false
    }
  }

  async function fetchBoard(id: string) {
    loading.value = true
    try {
      currentBoard.value = await boardService.get(id)
    } finally {
      loading.value = false
    }
  }

  async function createBoard(name: string, teamId: string) {
    const board = await boardService.create({ name, teamId })
    boards.value.push(board)
    return board
  }

  async function deleteBoard(id: string) {
    await boardService.delete(id)
    boards.value = boards.value.filter((b) => b.id !== id)
    if (currentBoard.value?.id === id) {
      currentBoard.value = null
    }
  }

  return { boards, currentBoard, loading, fetchBoards, fetchBoard, createBoard, deleteBoard }
})
