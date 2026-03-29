// @vitest-environment jsdom
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createVuetify } from 'vuetify'
import StatCard from './StatCard.vue'

const vuetify = createVuetify()

function mountStatCard(props: { title: string; value: number | string; icon: string; gradient: string; change?: number }) {
  return mount(StatCard, {
    props,
    global: {
      plugins: [vuetify],
    },
  })
}

describe('StatCard', () => {
  it('renders value and title', () => {
    // Arrange
    const props = {
      title: 'Total Tasks',
      value: 42,
      icon: 'mdi-check',
      gradient: 'linear-gradient(135deg, #6366f1, #a855f7)',
    }

    // Act
    const wrapper = mountStatCard(props)

    // Assert
    expect(wrapper.text()).toContain('Total Tasks')
    expect(wrapper.text()).toContain('42')
  })

  it('renders string value', () => {
    // Arrange
    const props = {
      title: 'Status',
      value: 'Active',
      icon: 'mdi-circle',
      gradient: 'linear-gradient(135deg, #10b981, #06b6d4)',
    }

    // Act
    const wrapper = mountStatCard(props)

    // Assert
    expect(wrapper.text()).toContain('Active')
  })

  it('shows change percentage when provided', () => {
    // Arrange
    const props = {
      title: 'Completed',
      value: 15,
      icon: 'mdi-check',
      gradient: 'linear-gradient(135deg, #6366f1, #a855f7)',
      change: 12,
    }

    // Act
    const wrapper = mountStatCard(props)

    // Assert
    expect(wrapper.text()).toContain('12%')
  })

  it('does not show change when not provided', () => {
    // Arrange
    const props = {
      title: 'Tasks',
      value: 5,
      icon: 'mdi-check',
      gradient: 'linear-gradient(135deg, #6366f1, #a855f7)',
    }

    // Act
    const wrapper = mountStatCard(props)

    // Assert
    expect(wrapper.text()).not.toContain('%')
  })
})
