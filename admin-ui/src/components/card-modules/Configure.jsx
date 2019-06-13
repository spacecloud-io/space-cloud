import React from 'react'

function Configure(props) {
  return (
    <div className="configure card-style">
      <i class="material-icons module">settings</i>
      <p className="heading">Configure</p>
      <div className="underline"></div>
      <span className="desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit,
        sed do eiusmod tempor incididunt ut. </span>
      <div className="footer">
        <a href="https://spaceuptech.com/docs/configure/overview" target="blank" >
          <span className="docs">Documentation</span>
          <i className="material-icons book">import_contacts</i>
        </a>
        {props.modules.enabled ?
          <button className="enabled-button">Enabled</button> :
          <button className="disabled-button">Disabled</button>
        }
      </div>
    </div>
  )
}

export default Configure

