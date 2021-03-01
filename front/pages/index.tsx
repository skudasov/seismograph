import React from 'react'
import ReactMarkdown from 'react-markdown'
import { Container } from '@material-ui/core'

const markdown = `
## Seismograph
Simple and reliable service for load tests, benchmarks.
`

export const Index = () => {
  return (
    <Container>
      <ReactMarkdown>{markdown}</ReactMarkdown>
    </Container>
  )
}

export default Index
