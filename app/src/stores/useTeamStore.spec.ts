// @vitest-environment jsdom
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useTeamStore } from './useTeamStore'

vi.mock('@/services/teamService', () => ({
  teamService: {
    list: vi.fn().mockResolvedValue([
      { id: '1', name: 'Engineering', organizationId: 'org-1', createdAt: '' },
      { id: '2', name: 'Design', organizationId: 'org-1', createdAt: '' },
    ]),
    get: vi.fn().mockResolvedValue({
      id: '1',
      name: 'Engineering',
      organizationId: 'org-1',
      createdAt: '',
      members: [
        {
          id: 'm1', teamId: '1', userId: 'u1', role: 'leader', joinedAt: '',
          user: { id: 'u1', name: 'John', email: 'john@test.com', avatarUrl: null, createdAt: '', updatedAt: '' },
        },
      ],
    }),
    create: vi.fn().mockResolvedValue({ id: '3', name: 'Marketing', organizationId: 'org-1', createdAt: '' }),
    delete: vi.fn().mockResolvedValue(undefined),
    invite: vi.fn().mockResolvedValue({
      id: 'm2', teamId: '1', userId: 'u2', role: 'member', joinedAt: '',
    }),
    removeMember: vi.fn().mockResolvedValue(undefined),
    updateMemberRole: vi.fn().mockResolvedValue(undefined),
  },
}))

describe('useTeamStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty teams', () => {
    const store = useTeamStore()
    expect(store.teams).toEqual([])
    expect(store.currentTeam).toBeNull()
  })

  it('fetchTeams loads teams', async () => {
    const store = useTeamStore()
    await store.fetchTeams()
    expect(store.teams).toHaveLength(2)
    expect(store.teams[0].name).toBe('Engineering')
  })

  it('fetchTeam loads team with members', async () => {
    const store = useTeamStore()
    await store.fetchTeam('1')
    expect(store.currentTeam).not.toBeNull()
    expect(store.currentTeam?.name).toBe('Engineering')
    expect(store.currentTeam?.members).toHaveLength(1)
    expect(store.currentTeam?.members[0].role).toBe('leader')
  })

  it('createTeam adds to list', async () => {
    const store = useTeamStore()
    await store.createTeam('Marketing', 'org-1')
    expect(store.teams).toHaveLength(1)
    expect(store.teams[0].name).toBe('Marketing')
  })

  it('deleteTeam removes from list', async () => {
    const store = useTeamStore()
    await store.fetchTeams()
    await store.deleteTeam('1')
    expect(store.teams).toHaveLength(1)
    expect(store.teams[0].id).toBe('2')
  })

  it('inviteMember adds to current team members', async () => {
    const store = useTeamStore()
    await store.fetchTeam('1')
    await store.inviteMember('1', 'jane@test.com', 'member')
    expect(store.currentTeam?.members).toHaveLength(2)
  })

  it('removeMember removes from current team', async () => {
    const store = useTeamStore()
    await store.fetchTeam('1')
    const initialLength = store.currentTeam?.members.length ?? 0
    await store.removeMember('1', 'u1')
    expect(store.currentTeam?.members.length).toBeLessThan(initialLength)
  })

  it('updateMemberRole changes role', async () => {
    const store = useTeamStore()
    await store.fetchTeam('1')
    await store.updateMemberRole('1', 'u1', 'member')
    expect(store.currentTeam?.members[0].role).toBe('member')
  })
})
