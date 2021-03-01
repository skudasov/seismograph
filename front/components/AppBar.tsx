import React from 'react'
import AppBar from '@material-ui/core/AppBar'
import Toolbar from '@material-ui/core/Toolbar'
import Button from '@material-ui/core/Button'
import IconButton from '@material-ui/core/IconButton'
import MenuIcon from '@material-ui/icons/Menu'
import HomeIcon from '@material-ui/icons/Home'
import Link from 'next/link'

export default function TopAppBar() {
  return (
    <div>
      <AppBar position="sticky">
        <Toolbar>
          <IconButton edge="start" color="inherit" aria-label="menu">
            <MenuIcon />
          </IconButton>
          <Link href="/">
            <IconButton edge="start" color="inherit" aria-label="menu">
              <HomeIcon />
            </IconButton>
          </Link>
          <Button color="inherit" href="/projects" data-test="menu-projects">
            Projects
          </Button>
          <Button color="inherit" href="/clusters" data-test="menu-projects">
            Clusters
          </Button>
        </Toolbar>
      </AppBar>
      <br />
    </div>
  )
}
