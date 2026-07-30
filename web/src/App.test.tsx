import { describe, it, expect } from 'vitest'

describe('App smoke', () => {
  it('renders heading', async () => {
    const mod = await import('./App')
    const App = mod.default
    const { container } = renderWithStrictMode(App)
    expect(container.innerHTML).toContain('GhostStack Dashboard')
  })
})
