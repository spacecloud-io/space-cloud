import React from 'react'
import './header.css'

function Header(props) {
  return (
    <div className="header" style={{ textAlign: props.align }}>
      <p className="heading" style={{ color: props.color, fontSize: props.fontSize, textAlign: props.align }}>{props.name}</p>
      <div className="line" style={{ textAlign: props.align }}></div>
    </div>
  )
}

export default Header
