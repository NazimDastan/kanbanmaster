// @vitest-environment jsdom
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { createVuetify } from 'vuetify'
import TaskCard from './TaskCard.vue'
import type { Task } from '@/types/task'

vi.mock('@/utils/date', () => ({
  formatDeadline: vi.fn((dateStr: string) => ({
    text: 'Dec 25',
    isOverdue: false,
    isUrgent: false,
  })),
}))

const vuetify = createVuetify()

function createTask(overrides: Partial<Task> = {}): Task {
  return {
    id: 't1',
    columnId: 'c1',
    title: 'Implement login page',
    description: null,
    creatorId: 'u1',
    assigneeId: null,
    assignee: null,
    priority: 'medium',
    deadline: null,
    position: 0,
    subtasks: [],
    labels: [],
    createdAt: '2026-01-01T00:00:00Z',
    updatedAt: '2026-01-01T00:00:00Z',
    completedAt: null,
    ...overrides,
  }
}

function mountTaskCard(task: Task) {
  return mount(TaskCard, {
    props: { task },
    global: {
      plugins: [vuetify],
    },
  })
}

describe('TaskCard', () => {
  it('renders task title', () => {
    // Arrange
    const task = createTask({ title: 'Build dashboard' })

    // Act
    const wrapper = mountTaskCard(task)

    // Assert
    expect(wrapper.text()).toContain('Build dashboard')
  })

  it('shows priority color stripe', () => {
    // Arrange
    const task = createTask({ priority: 'urgent' })

    // Act
    const wrapper = mountTaskCard(task)
    const stripe = wrapper.find('div.h-\\[2px\\]')

    // Assert
    expect(stripe.exists()).toBe(true)
    expect(stripe.attributes('style')).toContain('rgb(220, 38, 38)')
  })

  it('shows assignee avatar when assigned', () => {
    // Arrange
    const task = createTask({
      assigneeId: 'u2',
      assignee: {
        id: 'u2',
        name: 'Jane Doe',
        email: 'jane@example.com',
        avatarUrl: null,
        createdAt: '',
        updatedAt: '',
      },
    })

    // Act
    const wrapper = mountTaskCard(task)

    // Assert
    expect(wrapper.text()).toContain('JD')
  })

  it('does not show assignee avatar when unassigned', () => {
    // Arrange
    const task = createTask({ assignee: null })

    // Act
    const wrapper = mountTaskCard(task)
    const avatar = wrapper.find('.rounded-full')

    // Assert
    expect(avatar.exists()).toBe(false)
  })

  it('displays deadline text when deadline is set', () => {
    // Arrange
    const task = createTask({ deadline: '2026-12-25T00:00:00Z' })

    // Act
    const wrapper = mountTaskCard(task)

    // Assert
    expect(wrapper.text()).toContain('Dec 25')
  })

  it('emits click event with task when clicked', async () => {
    // Arrange
    const task = createTask()
    const wrapper = mountTaskCard(task)

    // Act
    await wrapper.find('div').trigger('click')

    // Assert
    expect(wrapper.emitted('click')).toBeTruthy()
    expect(wrapper.emitted('click')![0]).toEqual([task])
  })
})
