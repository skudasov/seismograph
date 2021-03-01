import React from 'react'
import ReactEcharts from 'echarts-for-react'

interface TickChartSeries {
  Tick: Array<string>
  90: Array<string>
  95: Array<string>
  99: Array<string>
}

interface TickChart {
  series: TickChartSeries
}

interface ChartTickPercentilesProps {
  chart: TickChart
}

const ChartTickPercentiles = (props: ChartTickPercentilesProps) => {
  const { chart } = props

  return (
    <div data-test="test-echart">
      <ReactEcharts
        option={{
          legend: {
            show: true,
            data: ['aa', 'bb'],
          },
          tooltip: {
            trigger: 'axis',
          },
          dataZoom: [
            {
              xAxisIndex: 0,
              minSpan: 5,
            },
          ],
          markLine: {
            data: [
              { type: 'average', name: 'trigger' },
              [
                {
                  symbol: 'none',
                  x: '90%',
                  yAxis: 'max',
                },
                {
                  symbol: 'circle',
                  label: {
                    position: 'start',
                    formatter: 'trigger',
                  },
                  type: 'max',
                  name: 'trigger',
                },
              ],
            ],
          },
          toolbox: {
            show: true,
            feature: {
              dataZoom: {
                yAxisIndex: 'none',
              },
              dataView: { readOnly: false },
              magicType: { type: ['line', 'bar'] },
              restore: {},
              saveAsImage: {},
            },
          },
          xAxis: {
            type: 'category',
            data: chart.series.Tick,
            name: 'Tick N',
          },
          yAxis: {
            type: 'value',
            name: 'Response time',
          },
          series: [
            {
              data: chart.series['90'],
              type: 'line',
            },
            {
              data: chart.series['95'],
              type: 'line',
            },
            {
              data: chart.series['99'],
              type: 'line',
            },
          ],
        }}
      />
    </div>
  )
}

export default ChartTickPercentiles
