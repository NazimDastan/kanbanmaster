import { describe, it, expect } from 'vitest'
import { truncate, pluralize, getInitials } from './format'

describe('truncate', () => {
  it('returns full text if under limit', () => {
    expect(truncate('hello', 10)).toBe('hello')
  })

  it('truncates long text', () => {
    expect(truncate('hello world', 5)).toBe('hello...')
  })
})

describe('pluralize', () => {
  it('returns singular for 1', () => {
    expect(pluralize(1, 'task')).toBe('task')
  })

  it('returns plural for > 1', () => {
    expect(pluralize(5, 'task')).toBe('tasks')
  })

  it('uses custom plural', () => {
    expect(pluralize(2, 'person', 'people')).toBe('people')
  })
})

describe('getInitials', () => {
  it('returns initials from full name', () => {
    expect(getInitials('John Doe')).toBe('JD')
  })

  it('limits to 2 characters', () => {
    expect(getInitials('John James Doe')).toBe('JJ')
  })

  it('handles single name', () => {
    expect(getInitials('John')).toBe('J')
  })
})
