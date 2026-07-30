import { describe, it, expect } from 'vitest'

describe('GhostStack SDK', () => {
  it('exports GhostStackClient class', async () => {
    const mod = await import('./ghoststack')
    expect(typeof mod.GhostStackClient).toBe('function')
  })
})
