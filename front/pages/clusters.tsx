import React, { useState } from 'react'
import { Container, Grid } from '@material-ui/core'
import {
  ColDef,
  DataGrid,
  RowParams,
  SelectionChangeParams,
} from '@material-ui/data-grid'
import { useRouter } from 'next/router'
import { useQuery } from 'react-query'
import styles from './project/[pid]/tests.module.css'
import NewClusterModal from '../components/modal/NewCluster'
import ClustersGridBar from '../components/ClustersGridBar'
import { apiFetchClusters } from '../api/clusters'
import { clustersToGrid } from '../data/functions/testsToGrid'

const ClustersMetaGrid = () => {
  const [popup, setPopup] = useState(false)

  const { isLoading, isError, data, error } = useQuery(
    'clusters',
    apiFetchClusters
  )

  const router = useRouter()

  const columns: ColDef[] = [
    { field: 'col1', headerName: 'Name', width: 200 },
    { field: 'col2', headerName: 'Provider name', width: 200 },
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
        <NewClusterModal popup={popup} setPopup={setPopup} />
        <Grid item>
          <ClustersGridBar setPopup={setPopup} />
        </Grid>
        <Grid item>
          <DataGrid
            autoHeight
            className={styles.dataGrid}
            rows={clustersToGrid(data)}
            columns={columns}
            onSelectionChange={rowSelected}
            onRowClick={(param: RowParams) => {
              router.push({
                pathname: `/attack_cluster/${param.rowModel.data.id}/vms`,
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

export default ClustersMetaGrid
