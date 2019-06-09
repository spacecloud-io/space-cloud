import React from 'react'
import logo from '../../assets/logo-black.svg';
import { connect } from 'react-redux';
import './topbar.css'
import { Button } from 'antd';

function Topbar(props){
 return (
     <div className="topbar">
       <img className="logo-black" src={logo} />
       <span>{props.title}</span>
       <Button type="primary" className="save-button" onClick={props.handleSave}>SAVE</Button>
     </div>
 )
}

const mapStateToProps = (state, ownProps) => {
  return {
    title: "User-management"
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleSave: () => {
      console.log('Saved')
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Topbar);