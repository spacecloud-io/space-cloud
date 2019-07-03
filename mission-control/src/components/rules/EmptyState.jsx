import React from 'react'
import { Button } from 'antd'
import './empty-state.css'

function EmptyState(props) {
  return (
    <div className="empty-state">
      <img src={props.graphics} alt="graphics" /><br />
      <p>{props.desc}</p>
      <Button type="primary" className="add-btn" onClick={props.handleClick}>{props.buttonText}</Button>
    </div>
  )
}

export default EmptyState
