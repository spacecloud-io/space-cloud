import React from 'react'
import mail from '../../assets/mail.svg'
import google from '../../assets/google.svg'
import fb from '../../assets/fb.svg'
import twitter from '../../assets/twitter.svg'
import github from '../../assets/github.svg'

function UserManagement(props) {
  return (
    <div className="user-management card-style">
      <div>
      <i className="material-icons module">people</i>
      <div className="logos">
      <img src={mail} alt="mail" />
      <img src={google} alt="google" />
      <img src={fb} alt="fb" />
      <img src={twitter} alt="twitter" />
      <img src={github} alt="github" />
      </div>
      </div>
      <p className="heading">User Management</p>
      <div className="underline"></div>
      <span className="desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit,
        sed do eiusmod tempor incididunt ut labore. </span><br />
      <div className="footer">
        <a href="https://spaceuptech.com/docs/user-management/overview" target="blank" >
          <span className="docs">Documentation</span>
          <i className="material-icons book">import_contacts</i>
        </a>
        <button className="button">Enabled</button>
        </div>
      </div>
      )
    }
    
    export default UserManagement
    
