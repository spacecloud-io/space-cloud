import React from 'react'
import { Link } from 'react-router-dom';
import { Divider } from "antd"
import SidenavItem from './SidenavItem'
import './sidenav.css'
import Header from '../header/Header'
import { connect } from 'react-redux'
import { get } from 'automate-redux';

const Sidenav = (props) => {
  return(
    <div className="sidenav">
    <div className="flex-container">
      <Header name={props.projectId} color="#000" fontSize="18px" />
    </div>
    <Link to={`/mission-control/projects/${props.projectId}/overview`}>
      <SidenavItem name="Project Overview" icon="home" active={props.selectedItem === 'overview'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/user-management`}>
      <SidenavItem name="User Management" icon="people" active={props.selectedItem === 'user-management'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/database`}>
      <SidenavItem name="Database" icon="dns" active={props.selectedItem === 'database'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/file-storage`}>
      <SidenavItem name="File Storage" icon="folder_open" active={props.selectedItem === 'file-storage'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/functions`}>
      <SidenavItem name="Functions" icon="code" active={props.selectedItem === 'functions'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/configure`}>
      <SidenavItem name="Configure" icon="settings" active={props.selectedItem === 'configure'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/explorer`}>
      <SidenavItem name="Explorer" icon="explore" active={props.selectedItem === 'explorer'} />
    </Link>
    <Divider />
    <Link to={`/mission-control/projects/${props.projectId}/gateway`}>
      <SidenavItem name="Gateway" icon="cloud" active={props.selectedItem === 'gateway'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/deploy`}>
      <SidenavItem name="Deploy" icon="local_airport" active={props.selectedItem === 'deploy'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/plans`}>
      <SidenavItem name="Plans" icon="assignment" active={props.selectedItem === 'plans'} />
    </Link>
    <Link to={`/mission-control/projects/${props.projectId}/billing`}>
      <SidenavItem name="Billing" icon="attach_money" active={props.selectedItem === 'billing'} />
    </Link>
  </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    projectId: get(state, "config.id", ""),
    selectedItem: ownProps.selectedItem,
  }
}


export default connect(mapStateToProps)(Sidenav);
