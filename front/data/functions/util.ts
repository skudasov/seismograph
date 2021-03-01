function randomString(length: number) {
  let result = ''
  const characters =
    'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
  const charactersLength = characters.length
  for (let i = 0; i < length; i++) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength))
  }
  return result
}

const randRange = (min: number, max: number): number => {
  return Math.floor(Math.random() * (max - min + 1) + min)
}

function randomPointCSV(
  header: string,
  pointsInLineNumber: number,
  pointsNumber: number
): string {
  const points = []
  points.push(header)
  const tick = 1
  for (let i = 0; i < pointsNumber; i++) {
    const row = []
    row.push(tick.toString())
    for (let j = 0; i < pointsInLineNumber; j++) {
      row.push(randRange(100, 1000).toString())
    }
    const jrow = row.join(',')
    points.push(jrow)
  }
  return points.join('\n')
}

export { randomString, randomPointCSV }
