import React from 'react'
import { Button, Col } from 'antd'
import './database-card.css'

function DatabaseCard(props) {
  return (
    <Col span={8} >
      <div className="db-card">
        <div className="db-image">
          <img src={props.graphics} alt={props.name} />
        </div>
        <div className="db-details">
          <div className="db-name"><p>{props.name}</p></div>
          <span>{props.desc}</span><br />
          <div className="db-enable"><Button type="primary" className="enable-btn" onClick={props.handleEnable}>Enable</Button></div>
        </div>
      </div>
    </Col>
  )
}

export default DatabaseCard