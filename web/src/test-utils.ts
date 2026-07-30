export async function renderWithStrictMode(ui: React.ReactElement) {
  const { default: testingLibrary } = await import('@testing-library/react')
  return testingLibrary.render(ui, { wrapper: undefined })
}
