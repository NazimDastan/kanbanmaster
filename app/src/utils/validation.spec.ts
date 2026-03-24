import { describe, it, expect } from 'vitest'
import { required, email, minLength, maxLength } from './validation'

describe('required', () => {
  it('returns error for empty string', () => {
    expect(required('Name')('')).toBe('Name is required')
  })

  it('returns true for non-empty string', () => {
    expect(required('Name')('John')).toBe(true)
  })
})

describe('email', () => {
  it('validates correct email', () => {
    expect(email('test@example.com')).toBe(true)
  })

  it('rejects invalid email', () => {
    expect(email('invalid')).toBe('Invalid email address')
  })

  it('rejects email without domain', () => {
    expect(email('test@')).toBe('Invalid email address')
  })
})

describe('minLength', () => {
  it('rejects short strings', () => {
    expect(minLength(6)('abc')).toBe('Must be at least 6 characters')
  })

  it('accepts valid length', () => {
    expect(minLength(6)('abcdef')).toBe(true)
  })
})

describe('maxLength', () => {
  it('rejects long strings', () => {
    expect(maxLength(5)('abcdef')).toBe('Must be less than 5 characters')
  })

  it('accepts valid length', () => {
    expect(maxLength(5)('abc')).toBe(true)
  })
})
