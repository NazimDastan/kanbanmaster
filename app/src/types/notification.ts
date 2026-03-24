export type NotificationType =
  | 'assigned'
  | 'delegated'
  | 'deadline'
  | 'comment'
  | 'report_request'
  | 'completed'
  | 'overdue'

export interface Notification {
  id: string
  userId: string
  type: NotificationType
  title: string
  message: string
  referenceId: string | null
  isRead: boolean
  createdAt: string
}
