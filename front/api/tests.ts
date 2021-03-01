import { randomString } from '../data/functions/util'

const apiFetchTests = async (options: any) => {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/project/${options.pid}/tests`
  )
  return res.json()
}

const apiFetchTest = async (options: any) => {
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/project/${options.pid}/test/${options.tid}`
  )
  return res.json()
}

const generateFakeTestData = async (pid: string | string[] | undefined) => {
  const formData = new FormData()
  const f = new Blob([
    'Tick,90,95,99\n1,100,200,300\n2,105,205,305\n3,211,280,325\n4,200,300,450\n5,300,306,600',
  ])

  formData.append('test_data_csv', f, 'test_data_csv')
  formData.append(
    'meta',
    JSON.stringify({ test_name: randomString(10), environment: 'dev' })
  )
  const res = await fetch(
    `${process.env.NEXT_PUBLIC_API_URL}/project/${pid}/test`,
    {
      method: 'POST',
      body: formData,
    }
  )
  return res.json()
}

export { apiFetchTests, apiFetchTest, generateFakeTestData }
