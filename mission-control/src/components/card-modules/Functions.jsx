import React from 'react'

function Functions(props) {
  return (
    <div className="overview-card functions-card">
      <i class="material-icons module">code</i>
      <p className="heading">Functions Mesh</p>
      <div className="underline"></div>
      <span className="desc">Architect your microservices in the form of functions! Invoke these functions from fronted or backend directly.</span>
      <div className="footer">
        <a href="https://spaceuptech.com/docs/functions/overview" target="blank" >
          <span className="docs">Documentation</span>
          <i className="material-icons book">import_contacts</i>
        </a>
        {/* {props.modules.enabled ?
          <button className="button">Enabled</button> :
          <button className="disabled button">Disabled</button>
        } */}
      </div>
    </div>
  )
}

export default Functions

