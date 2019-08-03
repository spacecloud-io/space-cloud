import React, { Component } from 'react'
import ReactGA from 'react-ga';
import { connect } from 'react-redux'
import { set, get } from "automate-redux"
import { cloneDeep } from "lodash"
import service from '../../index';
import store from "../../store"
import history from "../../history"
import { generateProjectConfig, notify, adjustConfig } from '../../utils';

import { Row, Col, Button, Form, Input, Icon } from 'antd'
import StarterTemplate from '../../components/starter-template/StarterTemplate'
import Topbar from '../../components/topbar/Topbar'
import './create-project.css'

import create from '../../assets/create.svg'
import postgresIcon from '../../assets/postgresIcon.svg'
import mysqlIcon from '../../assets/mysqlIcon.svg'
import mongoIcon from '../../assets/mongoIcon.svg'

class CreateProject extends Component {
  constructor(props) {
    super(props)
    this.state = {
      selected: "mongo"
    };
  }

  componentDidMount() {
		ReactGA.pageview("/create-project");
	}

  handleSelect(value) {
    return this.setState({ selected: value });
  }

  handleSubmit = e => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        this.props.handleNext(values.projectName, this.state.selected)
      }
    });
  };

  render() {
    const { getFieldDecorator } = this.props.form;
    return (
      <div className="create-project">
        <Topbar hideActions />
        <div className="content">
          <div>
            <span>PROJECT NAME</span>
            <Form>
              <Form.Item >
                {getFieldDecorator('projectName', {
                  rules: [{ required: true, message: 'Please input a project name' }],
                })(
                  <Input
                    prefix={<Icon type="edit" style={{ color: 'rgba(0,0,0,.25)' }} />}
                    placeholder="Project name" style={{ width: 600 }} />,
                )}
              </Form.Item>
            </Form>
          </div>

          <p>SELF-HOSTED</p>
          <div className="underline"></div>
          <div className="cards">
            <Row>
              <Col span={6}>
                <StarterTemplate icon={mongoIcon} onClick={() => this.handleSelect("mongo")}
                  heading="MONGODB" desc="A open-source cross-platform document- oriented database."
                  recommended={false} selected={this.state.selected}
                  active={this.state.selected === "mongo"} />
              </Col>

              <Col span={6}>
                <StarterTemplate icon={postgresIcon} onClick={() => this.handleSelect("sql-postgres")}
                  heading="POSTGRESQL" desc="The world's most advanced open source database."
                  recommended={false} selected={this.state.selected}
                  active={this.state.selected === "sql-postgres"} />
              </Col>

              <Col span={6}>
                <StarterTemplate icon={mysqlIcon} onClick={() => this.handleSelect("sql-mysql")}
                  heading="MYSQL" desc="The world's most popular open source database."
                  recommended={false} selected={this.state.selected}
                  active={this.state.selected === "sql-mysql"} />
              </Col>
            </Row>

          </div>
          <img className="image" src={create} alt="graphic" height="380" width="360" />
          <Button type="primary" htmlType="submit" className="next-btn" onClick={this.handleSubmit}>NEXT</Button>
        </div>
      </div>
    )
  }
}
const WrappedCreateProject = Form.create({})(CreateProject)

const mapDispatchToProps = (dispatch) => {
  return {
    handleNext: (name, dbType) => {
      const projectConfig = generateProjectConfig(name, dbType)
      service.saveProjectConfig(projectConfig).then(() => {
        const updatedProjects = [...get(store.getState(), "projects", []), projectConfig]
        dispatch(set("projects", updatedProjects))
        history.push(`/mission-control/projects/${projectConfig.id}`)
        const adjustedConfig = adjustConfig(projectConfig)
        dispatch(set("config", adjustedConfig))
        dispatch(set("savedConfig", cloneDeep(adjustedConfig)))
        notify("success", "Success", "Project created successfully with suitable defaults")
      }).catch(error => {
        console.log("Error", error)
        notify("error", "Error", "Could not create project")
      })
    }
  }
}

export default connect(null, mapDispatchToProps)(WrappedCreateProject);
