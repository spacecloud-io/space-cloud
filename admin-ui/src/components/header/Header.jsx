import React from 'react'
import './heading.css'

function Header(props) {
  return (
    <div>
      <p className="heading" style={{color: props.color, fontSize: props.fontSize}}>{props.name}</p>
      <div className="line"></div>
    </div>
  )
}

export default Header
