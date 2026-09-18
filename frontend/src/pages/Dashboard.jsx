import { useEffect, useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { api } from '../api.js'

export default function Dashboard() {
  const [polls, setPolls] = useState([])
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const user = JSON.parse(localStorage.getItem('user') || '{}')

  useEffect(() => {
    api.myPolls().then(setPolls).catch(e => setError(e.message))
  }, [])

  function logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    navigate('/')
  }

  const totalVotes = (poll) => poll.options.reduce((sum, o) => sum + o.votes, 0)

  return (
    <div className="container">
      <nav>
        <Link to="/dashboard">Dashboard</Link>
        <Link to="/create">Create Poll</Link>
        <button style={{ width: 'auto', float: 'right' }} onClick={logout}>Logout</button>
      </nav>
      <h2>Welcome {user.name || ''} 👋</h2>
      <Link to="/create"><button>+ Create Poll</button></Link>

      <h3 style={{ marginTop: 24 }}>My Polls</h3>
      {error && <p className="error">{error}</p>}
      {polls.length === 0 && <p>No polls yet — create your first one!</p>}
      {polls.map(poll => (
        <div key={poll.id} className="link-box" style={{ flexDirection: 'column', alignItems: 'flex-start' }}>
          <strong>{poll.question}</strong>
          <span>{totalVotes(poll)} votes</span>
          <Link to={`/poll/${poll.id}`}>View live results →</Link>
        </div>
      ))}
    </div>
  )
}
