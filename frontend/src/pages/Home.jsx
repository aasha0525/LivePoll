import { Link } from 'react-router-dom'

export default function Home() {
  return (
    <div className="container" style={{ textAlign: 'center' }}>
      <h1>🗳️ LivePoll</h1>
      <p>Create polls. Share. Vote. Watch results live.</p>
      <Link to="/login"><button>Login</button></Link>
      <Link to="/signup"><button>Sign Up</button></Link>
    </div>
  )
}
