import { Button } from '@material-ui/core'
import ArrowBackIosIcon from '@material-ui/icons/ArrowBackIos'
import React from 'react'

const BackToTestsLink = () => {
  return (
    <Button data-test="back-to-tests" color="inherit" href="../tests">
      <ArrowBackIosIcon />
    </Button>
  )
}

export default BackToTestsLink
