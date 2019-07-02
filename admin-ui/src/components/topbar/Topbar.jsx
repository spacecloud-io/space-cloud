import React from 'react'
import logo from '../../assets/logo-black.svg';
import { connect } from 'react-redux';
import './topbar.css'
import { Button } from 'antd';
import DbSelector from '../../components/db-selector/DbSelector'
import { isEqual } from "lodash"
import history from "../../history";
import store from "../../store";
import { get } from 'automate-redux';
import { saveConfig } from '../../actions/index';


function Topbar(props) {
  return (
    <div className="topbar">
      <img className="logo-black" src={logo} alt="logo" />
      <span>{props.title}</span>
      {(props.title === "Database") &&
        <DbSelector handleSelect={props.handleSelect} selectedDb={props.selectedDb} />
      }
      <Button type="primary" className="save-button" onClick={props.handleSave} disabled={!props.unsavedChanges}>SAVE</Button>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    selectedDb: ownProps.selectedDb,
    unsavedChanges: !isEqual(state.config, state.savedConfig)
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  const projectId = get(store.getState(), "config.id", "")
  return {
    handleSave: saveConfig,

    handleSelect(value) {
      history.push(`/mission-control/${projectId}/database/rules/${value}`)
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Topbar);