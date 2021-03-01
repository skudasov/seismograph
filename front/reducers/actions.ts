import { RowId } from '@material-ui/data-grid'

export const SET_SELECTED_TESTS = 'app/set-selected-test'
export const UNSET_SELECTED_TESTS = 'app/unset-selected-test'

export const setSelectedTestAction = (ids: RowId[]) => ({
  type: SET_SELECTED_TESTS,
  ids,
})
