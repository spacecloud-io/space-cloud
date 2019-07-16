import React from 'react'
import {Button} from 'antd'
import './upgrade.css'

function Upgrade() {
  return (
    <div className="upgrade">
      <p className="upgrade-desc">All the development features in Space Cloud are free forever! You can upgrade
        to unlock operational capabilities like easy deployment, reporting and metrics</p><br />
      <p className="billing-choose-plan">Choose one of our paid plans to upgrade</p><br/>
      <Button type="primary">Explore all plans</Button>
    </div>
  )
}

export default Upgrade
