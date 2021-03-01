import { CreateProjectDTO } from './dto'

const apiCreateProject = async (body: CreateProjectDTO) => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/project`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  })
  return res.json()
}

const apiFetchProjects = async () => {
  const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/projects`)
  return res.json()
}

export { apiFetchProjects, apiCreateProject }
