import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { api } from '../api.js'

export default function CreatePoll() {
  const [question, setQuestion] = useState('')
  const [options, setOptions] = useState(['', ''])
  const [error, setError] = useState('')
  const [created, setCreated] = useState(null)
  const navigate = useNavigate()

  function updateOption(i, value) {
    const next = [...options]
    next[i] = value
    setOptions(next)
  }

  function addOption() {
    setOptions([...options, ''])
  }

  function removeOption(i) {
    setOptions(options.filter((_, idx) => idx !== i))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    try {
      const poll = await api.createPoll({ question, options })
      setCreated(poll)
    } catch (err) {
      setError(err.message)
    }
  }

  if (created) {
    const link = `${window.location.origin}/poll/${created.id}`
    return (
      <div className="container">
        <h2>Poll Created! 🎉</h2>
        <p>Share this poll:</p>
        <div className="link-box">
          <span>{link}</span>
          <button style={{ width: 'auto' }} onClick={() => navigator.clipboard.writeText(link)}>Copy Link</button>
        </div>
        <Link to={`/poll/${created.id}`}><button>Watch live results</button></Link>
        <Link to="/dashboard"><button>Back to dashboard</button></Link>
      </div>
    )
  }

  return (
    <div className="container">
      <h2>Create Poll</h2>
      <form onSubmit={handleSubmit}>
        <label>Question</label>
        <input value={question} onChange={e => setQuestion(e.target.value)} placeholder="What is your favourite language?" required />

        <label>Options</label>
        {options.map((opt, i) => (
          <div className="option-row" key={i}>
            <input value={opt} onChange={e => updateOption(i, e.target.value)} placeholder={`Option ${i + 1}`} />
            {options.length > 2 && (
              <button type="button" style={{ width: 'auto' }} onClick={() => removeOption(i)}>✕</button>
            )}
          </div>
        ))}
        <button type="button" onClick={addOption}>+ Add Option</button>

        {error && <p className="error">{error}</p>}
        <button type="submit">Create Poll</button>
      </form>
    </div>
  )
}
