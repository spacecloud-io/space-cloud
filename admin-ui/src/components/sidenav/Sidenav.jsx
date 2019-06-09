import React, { Component } from 'react'
import SidenavItemList from './SidenavItemList'
import './sidenav.css'
import Header from '../header/Header'
import EditItemModal from '../edit-item-modal/EditItemModal'
import { connect } from 'react-redux';
import projectId from '../../assets/projectId.svg';

const items = [{ name: "Project Overview", key: "overview", icon: "home" },
{ name: "User management", key: "user-management", icon: "people" },
{ name: "Database", key: "database", icon: "dns" },
{ name: "Functions", key: "functions", icon: "code" },
{ name: "Configure", key: "configure", icon: "settings" }]

class Sidenav extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modalVisible: false
    };
    this.handleModalVisiblity=this.handleModalVisiblity.bind(this);
  }

  handleModalVisiblity(visible){
    this.setState({modalVisible: visible});
  }

  render() {
    return (
      <div className="sidenav">
        <div className="flex-container">
          <Header name={this.props.projectId} color="#000" fontSize="18px" />
          <button className="edit" onClick={() => this.handleModalVisiblity(true)}><b>Edit</b></button>
        </div>
        
        <SidenavItemList items={items} selectedItem={this.props.selectedItem} />
        <EditItemModal graphics={projectId} heading="Project Id" name="Give a project Id" desc="You need to use the same project ID to initialize the client." placeholder="Enter a project ID" initialValue={this.props.projectId} visible={this.state.modalVisible} handleCancel={() => this.handleModalVisiblity(false)} handleSubmit={this.props.handleSubmit}/>
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
