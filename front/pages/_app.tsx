import React from 'react'
import withRedux from 'next-redux-wrapper'
import { AppContext, AppInitialProps, AppProps } from 'next/app'
import Head from 'next/head'
import { Global, css } from '@emotion/core'
import { Provider } from 'react-redux'
import { Store } from 'redux'
import { ThemeProvider } from '@material-ui/styles'
import { createMuiTheme } from '@material-ui/core'
import { green, purple } from '@material-ui/core/colors'
import { QueryClient, QueryClientProvider } from 'react-query'
import { ReactQueryDevtools } from 'react-query/devtools'
import TopAppBar from '../components/AppBar'
import { configureStore } from '../store/configureStore'
import { normalize } from '../constants/normalize'

type Props = { store: Store } & AppInitialProps & AppProps

type AppPage<P = {}> = {
  (props: P): JSX.Element | null
  getInitialProps: ({ Component, ctx }: AppContext) => Promise<AppInitialProps>
}

const theme = createMuiTheme({
  palette: {
    primary: {
      main: purple[500],
    },
    secondary: {
      main: green[500],
    },
  },
})

const queryClient = new QueryClient()

const App: AppPage<Props> = ({ store, pageProps, Component }) => {
  return (
    <>
      <Head>
        <title>Seismograph</title>
      </Head>
      <QueryClientProvider client={queryClient}>
        <Provider store={store}>
          <Global
            styles={css`
              ${normalize}
            `}
          />
          <ThemeProvider theme={theme}>
            <TopAppBar />
            <Component {...pageProps} />
          </ThemeProvider>
        </Provider>
        <ReactQueryDevtools initialIsOpen={false} />
      </QueryClientProvider>
    </>
  )
}

App.getInitialProps = async ({ Component, ctx }: AppContext) => {
  return {
    pageProps: {
      ...(Component.getInitialProps
        ? await Component.getInitialProps(ctx)
        : {}),
    },
  }
}

export default withRedux(configureStore)(App)
