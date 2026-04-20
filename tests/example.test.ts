import { describe, it, expect } from 'vitest'

describe('example tests', () => {
  it('should pass basic assertion', () => {
    expect(1 + 1).toBe(2)
  })

  it('should handle string operations', () => {
    const result = 'hello'.toUpperCase()
    expect(result).toBe('HELLO')
  })
})
