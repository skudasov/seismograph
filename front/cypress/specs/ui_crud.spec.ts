import { randomString } from '../../data/functions/util'

describe('UI CRUD tests', () => {
  beforeEach(() => {
    cy.visit('/')
  })
  it('can create new project, add test, see example chart', () => {
    const projName = randomString(10)
    cy.getBySel('menu-projects').click()
    cy.getBySel('projects-bar-new-btn').click()
    cy.getBySel('projects-modal-name-inp').type(projName)
    cy.getBySel('projects-modal-create-btn').click()
    cy.getDataGridCellByName(projName).click()
    cy.getBySel('random-test-add').click()
    cy.reload()
    cy.getDataGridCellByPos(0, 1).click()
    cy.getBySel('test-echart')
    cy.getBySel('back-to-tests').click()
  })
})
