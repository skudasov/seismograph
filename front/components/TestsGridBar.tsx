import { Button, Hidden } from '@material-ui/core'
import React from 'react'
import { connect } from 'react-redux'
import { RowId } from '@material-ui/data-grid'
import { RootState } from '../reducers'

const mapState = (state: RootState) => ({
  testsSelection: state.app.testsSelection,
})

const connector = connect(mapState)

interface GridBarProps {
  testsSelection: RowId[]
}

const TestsGridBar = (props: GridBarProps) => {
  const { testsSelection } = props

  return (
    <Hidden lgUp={testsSelection.length < 2}>
      <Button>Compare</Button>
    </Hidden>
  )
}

export default connector(TestsGridBar)
