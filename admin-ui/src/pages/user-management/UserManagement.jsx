import React from 'react'
import { connect } from 'react-redux'
import { Collapse, Icon  } from 'antd';
import Oauth from '../../components/user-management/Oauth'
import Email from '../../components/user-management/Email'
import Sidenav from '../../components/sidenav/Sidenav'
import Topbar from '../../components/topbar/Topbar'
import Documentation from '../../components/documentation/Documentation'
import Header from '../../components/header/Header'
import mailIcon from '../../assets/mailIcon.svg'
import googleIcon from '../../assets/googleIcon.svg'
import fbIcon from '../../assets/fbIcon.svg'
import twitterIcon from '../../assets/twitterIcon.svg'
import githubIcon from '../../assets/githubIcon.svg'
import CollapseHeader from './CollapseHeader'

const Panel = Collapse.Panel;
function UserManagement(props) {
  return (
    <div className="user-management">
      <Topbar title="User-management" />
      <div className="flex-box">
        <Sidenav selectedItem="user-management" />
        <div className="page-content">
          <div className="header-flex">
            <Header name="Authentication" color="#000" fontSize="22px" />
            <Documentation url="https://spaceuptech.com/docs/user-management" />
          </div>
          <Collapse accordion expandIconPosition="right" expandIcon={({ isActive }) => <Icon type="right" rotate={isActive ? 270 : 90} />}>
            <Panel header={(<CollapseHeader icon={mailIcon} desc="Mail" />)}  key="1">
              <Email formState={props.email} />
            </Panel>
            <Panel header={(<CollapseHeader icon={googleIcon} desc="Google" />)} key="2">
              <Oauth formState={props.google} type= "Google" redirectUrl= "www.google.com" handleChange={(values) => props.handleChange("google", values)}/>
            </Panel >
            <Panel header={(<CollapseHeader icon={fbIcon} desc="Facebook" />)} key="3">
              <Oauth formState={props.fb} type= "Facebook" redirectUrl= "www.fb.com" handleChange={(values) => props.handleChange("fb", values)}/>
            </Panel>
            <Panel header={(<CollapseHeader icon={twitterIcon} desc="Twitter" />)} key="4">
              <Oauth formState={props.twitter} type= "Twitter" redirectUrl= "www.twitter.com" handleChange={(values) => props.handleChange("twitter", values)} />
            </Panel>
            <Panel header={(<CollapseHeader icon={githubIcon} desc="Github" />)} key="5">
              <Oauth formState={props.github} type= "Github" redirectUrl= "www.github.com" handleChange={(values) => props.handleChange("github", values)} />
            </Panel>
          </Collapse><br />
        </div> 
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    email: {enabled: true},
    google: {enabled: true, appId: "1045208832449477", appSecret: "adajdoidj121dicn32q"},
    fb: {enabled: true, appId: "2045208832449477", appSecret: "adajdoidj121dicn32q"},
    twitter: {enabled: true, appId: "3045208832449477", appSecret: "adajdoidj121dicn32q"},
    github: {enabled: true, appId: "4045208832449477", appSecret: "adajdoidj121dicn32q"}
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleChange: (provider, values) => {
      console.log(provider , values)
    },
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(UserManagement);

