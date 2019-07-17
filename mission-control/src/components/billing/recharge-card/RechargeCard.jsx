import React from "react"
import { Icon } from 'antd'
import './recharge-card.css'

function RechargeCard({ amount, extraCredits, handleClick, active, desc }) {
  return (
    <div className={`cards ${active ? 'selected' :  ''}`} onClick={handleClick}>
      <div className="pricing-cards">
        {active &&
          <div>
            <Icon type="star" theme="filled" className="star" />
            <div class="triangle-down"></div>
          </div>
        }

        <div className="cards-title">${amount}</div>
        <div className="cards-desc">{desc}</div>
        {extraCredits ?
          <div className="footer">Free {extraCredits}$</div> :
          <div className="footer">No free credits</div>}
      </div>
    </div>
  )
}

export default RechargeCard