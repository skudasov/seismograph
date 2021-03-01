import { CreateClusterDTO } from './dto'

const apiCreateCluster = async (body: CreateClusterDTO) => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/attack_cluster`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  })
  return res.json()
}

const apiFetchClusters = async () => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/attack_clusters`)
  return res.json()
}

export { apiFetchClusters, apiCreateCluster }
