import React from 'react'
import './header.css'

function Header(props) {
  return (
    <div className="header">
      <p className="heading" style={{color: props.color, fontSize: props.fontSize}}>{props.name}</p>
      <div className="line"></div>
    </div>
  )
}

export default Header
