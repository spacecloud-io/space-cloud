import React from 'react'

function SidenavItem(props) {
  return (
      <div className={
        props.active ? 'item active' : 'item'
      }>
      <i className="material-icons">{props.icon}</i>
      <span>{props.name}</span>
    </div>
  )
}

export default SidenavItem