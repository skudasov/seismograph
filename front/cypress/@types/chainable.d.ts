declare global {
  // eslint-disable-next-line no-redeclare
  namespace Cypress {
    interface Chainable<Subject> {
      getDataGridCellByPos(row: number, col: number): Chainable<Element>

      getDataGridCellByName(sel: string): Chainable<Element>

      getBySel(sel: string): Chainable<Element>

      getInputBySel(sel: string): Chainable<Element>

      getBySelLike(sel: string): Chainable<Element>
    }
  }
}
