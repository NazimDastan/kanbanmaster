type ValidationRule = (value: string) => true | string

export const required = (label: string): ValidationRule => {
  return (v: string) => !!v || `${label} is required`
}

export const email: ValidationRule = (v: string) => {
  const pattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return pattern.test(v) || 'Invalid email address'
}

export const minLength = (min: number): ValidationRule => {
  return (v: string) => v.length >= min || `Must be at least ${min} characters`
}

export const maxLength = (max: number): ValidationRule => {
  return (v: string) => v.length <= max || `Must be less than ${max} characters`
}
