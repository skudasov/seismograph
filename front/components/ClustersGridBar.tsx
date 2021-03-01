import { Button } from '@material-ui/core'
import React from 'react'

interface ClustersGridBarProps {
  setPopup: (flag: boolean) => void
}

const ClustersGridBar = (props: ClustersGridBarProps) => {
  const { setPopup } = props

  return (
    <Button
      variant="contained"
      color="primary"
      onClick={() => {
        setPopup(true)
      }}
      data-test="clusters-bar-new-btn"
    >
      New
    </Button>
  )
}

export default ClustersGridBar
