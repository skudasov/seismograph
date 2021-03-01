import { Button } from '@material-ui/core'
import React from 'react'

interface ProjectsGridBar {
  setPopup: (flag: boolean) => void
}

const ProjectsGridBar = (props: ProjectsGridBar) => {
  const { setPopup } = props

  return (
    <Button
      variant="contained"
      color="primary"
      onClick={() => {
        setPopup(true)
      }}
      data-test="projects-bar-new-btn"
    >
      New
    </Button>
  )
}

export default ProjectsGridBar
