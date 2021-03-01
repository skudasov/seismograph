import { RowData, RowsProp } from '@material-ui/data-grid'

const testsToGrid = (data: Array<object>): RowsProp => {
  const rows: RowsProp = []
  if (data === undefined || data.tests === undefined) {
    return rows
  }
  data.tests.forEach((d: any) => {
    const row: RowData = {
      id: d.id,
      col1: d.test_name,
      col2: d.environment,
    }
    rows.push(row)
  })
  return rows
}

const projectsToGrid = (data: Array<object>): RowsProp => {
  const rows: RowsProp = []
  if (data === undefined || data.projects === undefined) {
    return rows
  }
  data.projects.forEach((d: any) => {
    const row: RowData = {
      id: d.id,
      col1: d.name,
      col2: d.description,
      col3: d.repo_url,
    }
    rows.push(row)
  })
  return rows
}

const clustersToGrid = (data: Array<object>): RowsProp => {
  const rows: RowsProp = []
  if (data === undefined || data.clusters === undefined) {
    return rows
  }
  data.clusters.forEach((d: any) => {
    const row: RowData = {
      id: d.id,
      col1: d.name,
      col2: d.provider_name,
    }
    rows.push(row)
  })
  return rows
}

export { testsToGrid, projectsToGrid, clustersToGrid }
