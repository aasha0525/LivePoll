const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api'
export const WS_URL = import.meta.env.VITE_WS_URL || 'ws://localhost:8080/api/ws'

async function request(path, options = {}) {
  const token = localStorage.getItem('token')
  const headers = { 'Content-Type': 'application/json', ...(options.headers || {}) }
  if (token) headers.Authorization = `Bearer ${token}`

  const res = await fetch(`${API_URL}${path}`, { ...options, headers })
  const data = await res.json().catch(() => ({}))

  if (!res.ok) {
    throw new Error(data.error || 'Something went wrong')
  }
  return data
}

export const api = {
  signup: (body) => request('/signup', { method: 'POST', body: JSON.stringify(body) }),
  login: (body) => request('/login', { method: 'POST', body: JSON.stringify(body) }),
  createPoll: (body) => request('/polls', { method: 'POST', body: JSON.stringify(body) }),
  myPolls: () => request('/polls'),
  getPoll: (id) => request(`/polls/${id}`),
  vote: (id, optionIndex) =>
    request(`/polls/${id}/vote`, { method: 'POST', body: JSON.stringify({ option_index: optionIndex }) }),
}
