import React from 'react'

function Functions(props) {
  return (
    <div className="functions card-style">
      <i class="material-icons module">code</i>
      <p className="heading">Functions</p>
      <div className="underline"></div>
      <span className="desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit,
        sed do eiusmod tempor incididunt ut labore. </span>
      <div className="footer">
        <a href="https://spaceuptech.com/docs/functions/overview" target="blank" >
          <span className="docs">Documentation</span>
          <i className="material-icons book">import_contacts</i>
        </a>
        {props.modules.enabled ?
          <button className="button">Enabled</button> :
          <button className="disabled button">Disabled</button>
        }
      </div>
    </div>
  )
}

export default Functions

