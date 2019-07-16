import React from "react"

function RechargeCard({ amount, extraCredits, handleClick, active }) {
  return (
    <div onClick={handleClick}>
      <p>{active && "active"}{amount}$</p>
      {extraCredits && <p>Free {extraCredits}$</p>}
    </div>
  )
}

export default RechargeCard