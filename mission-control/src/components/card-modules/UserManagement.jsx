import React from 'react'
import mail from '../../assets/mail.svg'
import google from '../../assets/google.svg'
import fb from '../../assets/fb.svg'
import twitter from '../../assets/twitter.svg'
import github from '../../assets/github.svg'

function UserManagement(props) {
  return (
    <div className="overview-card user-management-card">
      <div>
        <i className="material-icons module">people</i>
        <div className="logos">
          {props.modules.email &&
            <img src={mail} alt="mail" height="20px" width="20px"/>
          }
          {props.modules.google &&
            <img src={google} alt="google" height="20px" width="20px"/>
          }
          {props.modules.fb &&
            <img src={fb} alt="fb" height="20px" width="20px"/>
          }
          {props.modules.twitter &&
            <img src={twitter} alt="twitter" height="20px" width="20px"/>
          }
          {props.modules.github &&
            <img src={github} alt="github" height="20px" width="20px"/>
          }
        </div>
      </div>
      <p className="heading">User Management</p>
      <div className="underline"></div>
      <span className="desc">Let your users signin to your app seamlessly through email and Oauth via the user management module of Space Cloud.  </span><br />
      <div className="footer">
        <a href="https://spaceuptech.com/docs/user-management/overview" target="blank" >
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

export default UserManagement

