import React from 'react'
import './plan.css'

function Plan({ name, desc, points, pricing, active, handleClick }) {
  return (
    <div className={`plan ${active ? 'active' : ''}`}  onClick={handleClick}>
      <p>{name}</p>
      <p>{desc}</p>
      <ul>
        {points.map(point => <li>{point}</li>)}
      </ul>
      <p>{pricing}</p>
    </div>
  )
}

export default Plan
