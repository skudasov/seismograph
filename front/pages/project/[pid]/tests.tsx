import React from 'react'
import { Container, Grid } from '@material-ui/core'
import { Dispatch } from 'redux'
import { ConnectedProps, connect } from 'react-redux'
import {
  ColDef,
  DataGrid,
  RowId,
  RowParams,
  SelectionChangeParams,
} from '@material-ui/data-grid'
import { useRouter } from 'next/router'
import { useQuery } from 'react-query'
import { RootState } from '../../../reducers'
import { setSelectedTestAction } from '../../../reducers/actions'
import SearchBar from '../../../components/search/Search'
import styles from './tests.module.css'

import GridBar from '../../../components/TestsGridBar'
import DebugBar from '../../../components/DebugBar'
import { testsToGrid } from '../../../data/functions/testsToGrid'
import { apiFetchTests } from '../../../api/tests'

const mapState = (state: RootState) => ({
  testsSelection: state.app.testsSelection,
})

const mapDispatch = (dispatch: Dispatch) => ({
  setSelectionId: (ids: RowId[]) => dispatch(setSelectedTestAction(ids)),
})

const connector = connect(mapState, mapDispatch)

type SearchBarProps = ConnectedProps<typeof connector>

const TestsMetaGrid = (props: SearchBarProps) => {
  const { setSelectionId } = props

  const router = useRouter()
  const { pid } = router.query

  const { isLoading, isError, data, error } = useQuery('tests', () => {
    return apiFetchTests({ pid })
  })

  const columns: ColDef[] = [
    { field: 'col1', headerName: 'Name', width: 200 },
    { field: 'col2', headerName: 'Envorinment', width: 150 },
  ]

  const rowSelected = (sel: SelectionChangeParams) => {
    setSelectionId(sel.rowIds)
  }

  return (
    <Container>
      <Grid
        container
        spacing={3}
        className={styles.dataGrid}
        justify="space-around"
        direction="column"
      >
        <Grid item>
          <DebugBar />
        </Grid>
        <Grid item>
          <SearchBar className={styles.searchBar} />
        </Grid>
        <Grid item>
          <GridBar />
        </Grid>
        <Grid item>
          <DataGrid
            autoHeight
            className={styles.dataGrid}
            rows={testsToGrid(data)}
            columns={columns}
            onSelectionChange={rowSelected}
            onRowClick={(param: RowParams) => {
              router.push({
                pathname: `/project/${pid}/test/${param.rowModel.data.id}`,
              })
            }}
            checkboxSelection
            disableSelectionOnClick
          />
        </Grid>
      </Grid>
    </Container>
  )
}

export default connector(TestsMetaGrid)
