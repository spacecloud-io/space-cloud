import React from "react"
import { Button } from "antd"

export default ({ graphics, desc, actionText, handleClick }) => {
  return (
    <div>
      <img src={graphics} /> <br />
      <p>{desc}</p>
      <Button type="primary" onClick={handleClick}>{actionText}</Button>
    </div>
  )
}