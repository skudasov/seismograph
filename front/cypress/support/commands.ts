// / <reference types="cypress" />

Cypress.Commands.add('getBySel', (selector: string, ...args) => {
  const resultSelector = selector
    .split(' ')
    .map((item) => `[data-test=${item}]`)
    .join(' ')

  return cy.get(resultSelector, ...args)
})

Cypress.Commands.add('getDataGridCellByName', (selector: string, ...args) => {
  return cy.get(`[data-value=${selector}]`, ...args)
})

Cypress.Commands.add(
  'getDataGridCellByPos',
  (row: number, col: number, ...args) => {
    return cy.get(`[data-rowindex=${row}][aria-colindex=${col}]`, ...args)
  }
)

Cypress.Commands.add('getInputBySel', (selector, ...args) => {
  return cy.get(
    `[data-test=${selector}] input, [data-test=${selector}] textarea[name]`,
    ...args
  )
})

Cypress.Commands.add('getBySelLike', (selector, ...args) => {
  return cy.get(`[data-test*=${selector}]`, ...args)
})
