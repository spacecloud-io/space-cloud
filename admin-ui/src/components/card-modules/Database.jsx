import React from 'react'
import postgresMono from '../../assets/postgresMono.svg'
import mongoMono from '../../assets/mongoMono.svg'
import mysqlMono from '../../assets/mysqlMono.svg'

function Database(props) {
  return (
    <div className="database card-style">
      <div>
        <i class="material-icons module">dns</i>
        <div className="logos">
          {props.modules.mysql &&
            <img src={mysqlMono} alt="mysqlMono.svg" height="26px" width="26px" />
          }
          {props.modules.postgres &&
            <img src={postgresMono} alt="postgreSQL" height="26px" width="26px" />
          }
          {props.modules.mongo &&
            <img src={mongoMono} alt="mongoMono.svg" height="26px" width="16px" />
          }
        </div>
      </div>
      <p className="heading">Database</p>
      <div className="underline"></div>
      <span className="desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit,
        sed do eiusmod tempor incididunt ut labore. </span>
      <div className="footer">
        <a href="https://spaceuptech.com/docs/database/overview" target="blank" >
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

export default Database

