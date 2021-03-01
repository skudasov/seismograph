import { Card, CardContent, CardHeader, Divider } from '@material-ui/core'
import React from 'react'
import styles from './ReportCard.module.css'

interface TestMetaInfoProps {
  title: string
  children: React.ReactNode
}

const ReportCard = (props: TestMetaInfoProps) => {
  const { title, children } = props
  return (
    <Card className={styles.card} variant="outlined">
      <CardHeader title={title} />
      <Divider />
      <CardContent className={styles.card}>{children}</CardContent>
    </Card>
  )
}

export default ReportCard
