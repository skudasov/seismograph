import { BottomNavigation, BottomNavigationAction } from '@material-ui/core'
import React from 'react'
import { useRouter } from 'next/router'
import { generateFakeTestData } from '../api/tests'

const DebugBar = () => {
  const router = useRouter()

  const { pid } = router.query

  return (
    <BottomNavigation showLabels>
      <BottomNavigationAction
        data-test="random-test-add"
        label="Add random test"
        onClick={() => {
          generateFakeTestData(pid)
        }}
      />
    </BottomNavigation>
  )
}

export default DebugBar
