import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api, WS_URL } from '../api.js'

export default function PollView() {
  const { id } = useParams()
  const [poll, setPoll] = useState(null)
  const [error, setError] = useState('')
  const [voted, setVoted] = useState(false)
  const [selected, setSelected] = useState(null)
  const wsRef = useRef(null)

  // Load the poll once, then open a WebSocket that streams live vote updates.
  useEffect(() => {
    api.getPoll(id).then(setPoll).catch(e => setError(e.message))

    const ws = new WebSocket(`${WS_URL}/polls/${id}`)
    wsRef.current = ws

    ws.onmessage = (event) => {
      const update = JSON.parse(event.data) // { poll_id, option_index, votes }
      setPoll(prev => {
        if (!prev) return prev
        const options = [...prev.options]
        options[update.option_index] = { ...options[update.option_index], votes: update.votes }
        return { ...prev, options }
      })
    }

    return () => ws.close()
  }, [id])

  async function handleVote() {
    if (selected === null) return
    try {
      await api.vote(id, selected)
      setVoted(true)
    } catch (err) {
      setError(err.message)
    }
  }

  if (error) return <div className="container"><p className="error">{error}</p></div>
  if (!poll) return <div className="container"><p>Loading...</p></div>

  const total = poll.options.reduce((sum, o) => sum + o.votes, 0)

  return (
    <div className="container">
      <span className="live-badge">🟢 LIVE</span>
      <h2>{poll.question}</h2>

      {!voted && (
        <div>
          {poll.options.map((opt, i) => (
            <label key={i} style={{ display: 'block', margin: '8px 0' }}>
              <input type="radio" name="option" checked={selected === i} onChange={() => setSelected(i)} style={{ width: 'auto', marginRight: 8 }} />
              {opt.text}
            </label>
          ))}
          <button onClick={handleVote} disabled={selected === null}>Vote</button>
        </div>
      )}

      <div style={{ marginTop: 20 }}>
        {poll.options.map((opt, i) => {
          const pct = total > 0 ? Math.round((opt.votes / total) * 100) : 0
          return (
            <div className="bar-row" key={i}>
              <div className="bar-label"><span>{opt.text}</span><span>{opt.votes}</span></div>
              <div className="bar-bg"><div className="bar-fill" style={{ width: `${pct}%` }} /></div>
            </div>
          )
        })}
        <p>Total votes: {total}</p>
      </div>
    </div>
  )
}
