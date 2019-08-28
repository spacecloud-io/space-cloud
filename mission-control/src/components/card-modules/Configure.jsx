import React from 'react'

function Configure(props) {
  return (
    <div className="overview-card configure-card">
      <i class="material-icons module">settings</i>
      <p className="heading">Configure</p>
      <div className="underline"></div>
      <span className="desc">Configure all the modules of Space Cloud and how they plug into the system. Take complete control over your backend.</span>
      <div className="footer">
        <a href="https://spaceuptech.com/docs/" target="blank" >
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

export default Configure

