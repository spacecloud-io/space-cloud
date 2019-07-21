import React from "react"
import { Button } from "antd"

export default ({ graphics, desc, actionText, handleClick }) => {
  return (
    <div className="empty-state">
      <img src={graphics} /> <br />
      <p className="desc">{desc}</p>
      <Button type="primary" onClick={handleClick}>{actionText}</Button>
    </div>
  )
}