import React from 'react'
import { useRouter } from 'next/router'
import { Container, Grid, Typography } from '@material-ui/core'
import { useQuery } from 'react-query'
import ChartTickPercentiles from '../../../../components/charts/tickPercentiles'

import DefaultCard from '../../../../components/ReportCard'
import BackToTestsLink from '../../../../components/BackToTestsLink'
import { apiFetchTest } from '../../../../api/tests'

const TestView = () => {
  const router = useRouter()
  const { pid, tid } = router.query

  const { isLoading, isError, data, error } = useQuery('test', () => {
    return apiFetchTest({ pid, tid })
  })

  return (
    <div>
      <Container>
        <BackToTestsLink />
        <Grid container justify="space-between" direction="row">
          <Grid item xl={6}>
            {data && data.chart ? (
              <DefaultCard data-test="test-info" title="Info">
                <Typography color="textPrimary" gutterBottom>
                  Test: #{tid}
                </Typography>
                <Typography color="textPrimary" gutterBottom>
                  Store path: {data.chart.data_blob_path}
                </Typography>
                <Typography color="textPrimary" gutterBottom>
                  Name: {data.chart.test_name}
                </Typography>
                <Typography color="textPrimary" gutterBottom>
                  Environment: {data.chart.environment}
                </Typography>
              </DefaultCard>
            ) : (
              <div>Error</div>
            )}
          </Grid>
          <Grid item xl={6}>
            <DefaultCard title="Parameters">
              <Typography color="textPrimary" gutterBottom>
                VU: 100
              </Typography>
              <Typography color="textPrimary" gutterBottom>
                Test time: 12m28s
              </Typography>
              <Typography color="textPrimary" gutterBottom>
                <div>
                  <a href="https://github.com/insolar/soveren-integrations-sdk/blob/master/data/csv.go">
                    Test url
                  </a>
                </div>
              </Typography>
              <Typography color="textPrimary" gutterBottom>
                Retries: 3
              </Typography>
            </DefaultCard>
          </Grid>
          <Grid item xl={12}>
            <DefaultCard title="Response time percentiles">
              {data && data.chart ? (
                <ChartTickPercentiles chart={data.chart} />
              ) : (
                <div>Error</div>
              )}
            </DefaultCard>
          </Grid>
        </Grid>
      </Container>
    </div>
  )
}

export default TestView
