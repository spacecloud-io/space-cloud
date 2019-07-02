import React from 'react'
import './documentation.css'

function Documentation(props) {
  return (
    <div className="docs">
      <a href={props.url} target="blank" >
        Documentation
      <i className="material-icons book">import_contacts</i>
      </a>
    </div>
  )
}

export default Documentation