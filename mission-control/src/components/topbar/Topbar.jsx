import React, { Component } from 'react'
import { connect } from 'react-redux';
import { set, get } from 'automate-redux';
import { isEqual } from "lodash"
import service from "../../index"
import history from "../../history";
import store from "../../store";
import { unAdjustConfig, notify, openPlansPage } from "../../utils"

import { Button, Icon } from 'antd';
import DbSelector from '../../components/db-selector/DbSelector'
import SelectProject from '../../components/select-project/SelectProject'
import './topbar.css'

import logo from '../../assets/logo-black.svg';

class Topbar extends Component {
  constructor(props) {
    super(props)
    this.state = {
      modalVisible: false
    };
  }

  handleModalVisible(visible) {
    this.setState({ modalVisible: visible });
  }

  render() {
    return (
      <div>
        <div className="topbar">
          <img className="logo-black" src={logo} alt="logo" />
          {this.props.showProjectSelector && <div className="btn-position">
            <Button className="btn" onClick={() => this.handleModalVisible(true)}>{this.props.projectName}
              <Icon type="caret-down" /></Button>
          </div>}
          {(this.props.showDbSelector) &&
            <DbSelector handleSelect={this.props.handleSelect} selectedDb={this.props.selectedDb} />
          }
          {
            !this.props.hideActions && <div className="right-list">
              {this.props.mode < 1 && <Button type="primary" className="action-button upgrade-button" onClick={openPlansPage}>UPGRADE</Button>}
              <Button type="primary" className="action-button save-button" onClick={this.props.handleSave} disabled={!this.props.unsavedChanges}>SAVE</Button>
            </div>
          }
          <SelectProject visible={this.state.modalVisible} handleCancel={() => this.handleModalVisible(false)} />
        </div>
      </div>
    )
  }
}


const mapStateToProps = (state, ownProps) => {
  return {
    mode: get(state, "operationConfig.mode", 0),
    selectedDb: ownProps.selectedDb,
    projectName: get(state, "config.name", ""),
    unsavedChanges: !isEqual(state.config, state.savedConfig) || !isEqual(state.deployConfig, state.savedDeployConfig) || !isEqual(state.staticConfig, state.savedStaticConfig),
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    handleSave: () => {
      const config = get(store.getState(), "config")
      const deployConfig = get(store.getState(), "deployConfig")
      const staticConfig = get(store.getState(), "staticConfig")
      const mode = get(store.getState(), "operationConfig.mode", 0)

      // UnAdjust the config and check if there are any errors
      const result = unAdjustConfig(config, staticConfig)
      if (!result.ack) {
        if (Object.keys(result.errors.crud).length) {
          let errorDesc = ''
          Object.keys(result.errors.crud).forEach(db => {
            errorDesc = `${errorDesc} ${db} - (${result.errors.crud[db].join(", ")})`
          })
          notify("error", 'Error in CRUD config', errorDesc, 0)
        }

        if (result.errors.functions.length) {
          notify("error", 'Error in Functions config', `Services - ${result.errors.functions.join(", ")}`, 0)
        }

        if (result.errors.fileStore.length) {
          notify("error", 'Error in File Storage config', `Rules - ${result.errors.fileStore.join(", ")}`, 0)
        }

        if (result.errors.static.length) {
          notify("error", 'Error in Static Module config', `Rules - ${result.errors.static.join(", ")}`, 0)
        }
        return
      }

      dispatch(set("pendingRequests", true))
      const promises = mode > 0 ?
        [service.saveProjectConfig(result.config), service.saveStaticConfig(result.staticConfig), service.saveDeployConfig(deployConfig)]
        : [service.saveProjectConfig(result.config), service.saveStaticConfig(result.staticConfig)]
      Promise.all(promises)
        .then(() => {
          notify("success", 'Success', 'Config saved successfully')
          dispatch(set("savedConfig", config))
          dispatch(set("savedStaticConfig", staticConfig))
          if (mode > 0) {
            dispatch(set("savedDeployConfig", deployConfig))
          }
        })
        .catch(error => {
          console.log("Error", error)
          notify("error", "Error", 'Could not save config')
        })
        .finally(() => dispatch(set("pendingRequests", false)))
    },

    handleSelect(value) {
      const projectId = get(store.getState(), "config.id", "")
      history.push(`/mission-control/projects/${projectId}/database/rules/${value}`)
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Topbar);