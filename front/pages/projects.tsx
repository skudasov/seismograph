import React, { useState } from 'react'
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
import { RootState } from '../reducers'
import { setSelectedTestAction } from '../reducers/actions'
import styles from './project/[pid]/tests.module.css'
import ProjectsGridBar from '../components/ProjectsGridBar'
import NewProjectModal from '../components/modal/NewProject'
import { apiFetchProjects } from '../api/projects'
import { projectsToGrid } from '../data/functions/testsToGrid'

const mapState = (state: RootState) => ({
  projectsSelection: state.app.projectsSelection,
})

const mapDispatch = (dispatch: Dispatch) => ({
  setSelectionId: (ids: RowId[]) => dispatch(setSelectedTestAction(ids)),
})

const connector = connect(mapState, mapDispatch)

type Props = ConnectedProps<typeof connector>

const ProjectsMetaGrid = (props: Props) => {
  const { setSelectionId } = props

  const [popup, setPopup] = useState(false)

  const router = useRouter()

  const columns: ColDef[] = [
    { field: 'col1', headerName: 'Name', width: 200 },
    { field: 'col2', headerName: 'Description', width: 200 },
    { field: 'col3', headerName: 'Repo URL', width: 200 },
  ]

  const rowSelected = (sel: SelectionChangeParams) => {
    setSelectionId(sel.rowIds)
  }

  const { isLoading, isError, data, error } = useQuery(
    'projects',
    apiFetchProjects
  )

  return (
    <Container>
      <Grid
        container
        spacing={3}
        className={styles.dataGrid}
        justify="space-around"
        direction="column"
      >
        <NewProjectModal popup={popup} setPopup={setPopup} />
        <Grid item>
          <ProjectsGridBar setPopup={setPopup} />
        </Grid>
        <Grid item>
          <DataGrid
            autoHeight
            className={styles.dataGrid}
            rows={projectsToGrid(data)}
            columns={columns}
            onSelectionChange={rowSelected}
            onRowClick={(param: RowParams) => {
              router.push({
                pathname: `/project/${param.rowModel.data.id}/tests`,
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

export default connector(ProjectsMetaGrid)
