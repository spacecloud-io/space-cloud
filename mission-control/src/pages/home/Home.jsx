import React from "react"
import "./home.css"
import SpaceUpLogo from "../../logo.png"

export default () => {
  return (
    <div className="home">
      <div className="content">
        <span className="title">Welcome to Mission Control!</span>
        <br />
        <img src={SpaceUpLogo} alt="" />
        <br />
        <span>Loading data...</span>
      </div>
    </div>
  )
}