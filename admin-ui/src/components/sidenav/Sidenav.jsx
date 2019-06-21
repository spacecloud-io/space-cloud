import React, { Component } from 'react'
import { Link } from 'react-router-dom';

import SidenavItem from './SidenavItem'
import './sidenav.css'
import Header from '../header/Header'
import EditItemModal from '../edit-item-modal/EditItemModal'
import { connect } from 'react-redux'
import projectId from '../../assets/projectId.svg'

class Sidenav extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modalVisible: false
    };
    this.handleModalVisiblity = this.handleModalVisiblity.bind(this);
  }

  handleModalVisiblity(visible) {
    this.setState({ modalVisible: visible });
  }

  render() {
    return (
      <div className="sidenav">
        <div className="flex-container">
          <Header name={this.props.projectId} color="#000" fontSize="18px" />
          <button className="edit" onClick={() => this.handleModalVisiblity(true)}><b>Edit</b></button>
        </div>
        <Link to={`/mission-control/${this.props.projectId}/overview`}>
          <SidenavItem name="Project Overview" icon="home" active={this.props.selectedItem === 'overview'} />
        </Link>
        <Link to={`/mission-control/${this.props.projectId}/user-management`}>
          <SidenavItem name="User Management" icon="people" active={this.props.selectedItem === 'user-management'} />
        </Link>
        <Link to={`/mission-control/${this.props.projectId}/database`}>
          <SidenavItem name="Database" icon="dns" active={this.props.selectedItem === 'database'} />
        </Link>
        <Link to={`/mission-control/${this.props.projectId}/file-storage`}>
          <SidenavItem name="File Storage" icon="folder_open" active={this.props.selectedItem === 'file-storage'} />
        </Link>
        <Link to={`/mission-control/${this.props.projectId}/functions`}>
          <SidenavItem name="Functions" icon="code" active={this.props.selectedItem === 'functions'} />
        </Link>
        <Link to={`/mission-control/${this.props.projectId}/configure`}>
          <SidenavItem name="Configure" icon="settings" active={this.props.selectedItem === 'configure'} />
        </Link>
        <EditItemModal graphics={projectId} heading="Project ID" name="Give a project ID" desc="You need to use the same project ID to initialize the client." placeholder="Enter a project ID" initialValue={this.props.projectId} visible={this.state.modalVisible} handleCancel={() => this.handleModalVisiblity(false)} handleSubmit={this.props.handleSubmit} />
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => {
  return {
    projectId: "Todo-App",
    selectedItem: ownProps.selectedItem,
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleSubmit: (projectId) => {
      console.log('Submitted:', projectId)
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Sidenav);
