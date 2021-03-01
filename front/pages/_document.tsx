import Document, { Head, Html, Main, NextScript } from 'next/document'
import React from 'react'

interface Props {
  styleTags: React.ReactElement<{}>[]
}

export default class MyDocument extends Document<Props> {
  public render() {
    return (
      <Html>
        <Head />
        <body>
          <Main />
          <NextScript />
        </body>
      </Html>
    )
  }
}
