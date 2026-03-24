import type { Task } from './task'

export interface Board {
  id: string
  name: string
  teamId: string
  createdAt: string
}

export interface Column {
  id: string
  boardId: string
  name: string
  position: number
  tasks: Task[]
  createdAt: string
}

export interface BoardWithColumns extends Board {
  columns: Column[]
}
