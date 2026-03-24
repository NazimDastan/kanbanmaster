import type { User, UserRole } from './user'

export interface Organization {
  id: string
  name: string
  ownerId: string
  createdAt: string
}

export interface Team {
  id: string
  name: string
  organizationId: string
  createdAt: string
}

export interface TeamMember {
  id: string
  teamId: string
  userId: string
  user: User
  role: UserRole
  joinedAt: string
}

export interface TeamWithMembers extends Team {
  members: TeamMember[]
}
