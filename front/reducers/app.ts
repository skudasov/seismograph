import { SET_SELECTED_TESTS } from './actions'

export const initialState = {
  title: '',
  testsSelection: [],
  projectsSelection: [],
  newProjectModal: false,
  clustersSelection: [],
  newClusterModal: false,
}

export interface AppState {
  title: string
  testsSelection: string[]
  projectsSelection: string[]
  newProjectModal: boolean
  clustersSelection: string[]
  newClusterModal: boolean
}

export const reducer = (state = initialState, action: any) => {
  switch (action.type) {
    case SET_SELECTED_TESTS:
      return {
        ...state,
        testsSelection: action.ids,
      }
    default:
      return state
  }
}
