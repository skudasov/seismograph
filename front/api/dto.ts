export interface CreateProjectDTO {
  name: string
  description: string
  repo_url: string
}

export interface CreateClusterDTO {
  name: string
  provider_name: string
  project_id: number
  cluster_spec: ClusterSpecDTO
}

export interface ClusterSpecDTO {
  region: string
  instances: InstanceDTO[]
}

export interface InstanceDTO {
  region: string
  name: string
  image: string
  type: string
}
