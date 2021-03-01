module.exports = {
  target: process.env.NODE_ENV === 'production' ? 'serverless' : 'server',
  typescript: {
    ignoreBuildErrors: true,
  },
}
