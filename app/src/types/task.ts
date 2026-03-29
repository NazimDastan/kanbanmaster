import type { User } from './user'

export type Priority = 'urgent' | 'high' | 'medium' | 'low'

export interface Task {
  id: string
  columnId: string
  title: string
  description: string | null
  creatorId: string
  assigneeId: string | null
  assignee: User | null
  assignees: User[]
  priority: Priority
  deadline: string | null
  position: number
  subtasks: Subtask[]
  labels: Label[]
  createdAt: string
  updatedAt: string
  completedAt: string | null
}

export interface Subtask {
  id: string
  taskId: string
  title: string
  isCompleted: boolean
  createdAt: string
}

export interface Label {
  id: string
  boardId: string
  name: string
  color: string
}

export interface TaskDelegation {
  id: string
  taskId: string
  fromUserId: string
  toUserId: string
  reason: string
  delegatedAt: string
}

export interface Comment {
  id: string
  taskId: string
  userId: string
  user: User
  content: string
  createdAt: string
}

export interface ActivityLog {
  id: string
  taskId: string
  userId: string
  user: User
  action: string
  details: Record<string, unknown>
  createdAt: string
}

export const PRIORITY_CONFIG: Record<Priority, { label: string; color: string; icon: string }> = {
  urgent: { label: 'Urgent', color: '#dc2626', icon: 'mdi-alert-circle' },
  high: { label: 'High', color: '#06b6d4', icon: 'mdi-arrow-up-bold' },
  medium: { label: 'Medium', color: '#6366f1', icon: 'mdi-minus' },
  low: { label: 'Low', color: '#64748b', icon: 'mdi-arrow-down-bold' },
}
