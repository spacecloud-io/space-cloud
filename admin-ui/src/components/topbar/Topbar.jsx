import React from 'react'
import logo from '../../assets/logo-black.svg';
import { connect } from 'react-redux';
import './topbar.css'
import { Button } from 'antd';
import DbSelector from '../../components/db-selector/DbSelector'


function Topbar(props) {
  return (
    <div className="topbar">
      <img className="logo-black" src={logo} alt="logo" />
      <span>{props.title}</span>
      {props.title === "Database" &&
        <DbSelector handleSelect={props.handleSelect} selectedDb={props.selectedDb} />
      }
      <Button type="primary" className="save-button" onClick={props.handleSave}>SAVE</Button>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    selectedDb: 'sql-mysql',
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleSave: () => {
      console.log('Saved')
    },

    handleSelect(value) {
      console.log(`selected ${value}`);
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Topbar);