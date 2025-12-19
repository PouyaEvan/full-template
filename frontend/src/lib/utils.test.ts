import { cn } from './utils'

describe('cn', () => {
  test('merges class names correctly', () => {
    expect(cn('c-1', 'c-2')).toBe('c-1 c-2')
  })

  test('handles conditional classes', () => {
    expect(cn('c-1', true && 'c-2', false && 'c-3')).toBe('c-1 c-2')
  })

  test('merges tailwind classes', () => {
    expect(cn('px-2 py-1', 'p-4')).toBe('p-4')
  })
})
