import React from 'react'
import { Icon, Row, Col } from "antd"
import Header from "../../components/header/Header"
import './plan.css'

function Plan({ name, desc, points, pricing, active, handleClick }) {
  return (
    <div onClick={handleClick} className={`plan-card ${active ? 'selected' : ''}`}>
      <div className="content">
        <Header name={name} color="#ff6700" fontSize="18px" align="center" />
        <p className="center">{desc}</p>
        <Row>
          {points.map(point =>
            <Col span={20} offset={2}>
              <span>
                <Icon type="right" />
              </span>
              <span className="pointer">
                Realtime Crud API for any database
              </span>
            </Col>
          )}
        </Row>
        {active && <div className="tag-container center">
          <div className="tag">
            Active
        </div>
        </div>}
        <ul>
        </ul>

      </div>
      <div className="footer center">{pricing}</div>
    </div>
  )
}

export default Plan
