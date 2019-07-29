import React, { useState } from 'react'
import ReactGA from 'react-ga';
import { Link } from "react-router-dom"
import './welcome.css'
import { Button } from 'antd'

function Welcome(props) {
  useState(() => {
    ReactGA.pageview("/welcome");
  }, [])
  return (
    <div className="welcome-page">
      <div className="outer-rectangle"></div>
      <div className="inner-rectangle">
        <div className="content">
          <span className="welcome">Welcome to Space Cloud!</span>
          <div className="text">Develop great applications without having to write backend code.
          Focus more on business and less on technology.</div>
          <Link to="/mission-control/create-project"><Button type="primary" className="create-btn">CREATE A PROJECT</Button></Link>
        </div>
      </div>
    </div>
  )
}

export default Welcome;
