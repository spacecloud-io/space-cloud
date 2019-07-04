import React, { Component } from 'react'
import logo from '../../assets/logo-black.svg';
import { connect } from 'react-redux';
import './topbar.css'
import { Button, Icon } from 'antd';
import DbSelector from '../../components/db-selector/DbSelector'
import { isEqual } from "lodash"
import history from "../../history";
import store from "../../store";
import { get } from 'automate-redux';
import { saveConfig } from '../../actions/index';
import SelectProject from '../../components/select-project/SelectProject'

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
          {(this.props.save !== "false") &&
            <Button type="primary" className="save-button" onClick={this.props.handleSave} disabled={!this.props.unsavedChanges}>SAVE</Button>
          }
          <SelectProject visible={this.state.modalVisible} handleCancel={() => this.handleModalVisible(false)} />
        </div>
      </div>
    )
  }
}


const mapStateToProps = (state, ownProps) => {
  return {
    selectedDb: ownProps.selectedDb,
    projectName: get(state, "config.name", ""),
    unsavedChanges: !isEqual(state.config, state.savedConfig),
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    handleSave: saveConfig,

    handleSelect(value) {
      const projectId = get(store.getState(), "config.id", "")
      history.push(`/mission-control/projects/${projectId}/database/rules/${value}`)
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Topbar);