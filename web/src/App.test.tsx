import { describe, it, expect } from 'vitest'
import { renderWithStrictMode } from './test-utils'
import App from './App'

describe('App smoke', () => {
  it('renders heading', async () => {
    const result = await renderWithStrictMode(<App />)
    expect(result.container.innerHTML).toContain('GhostStack Dashboard')
  })
})
