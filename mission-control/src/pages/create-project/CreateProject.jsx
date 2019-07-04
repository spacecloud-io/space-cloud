import React, { Component } from 'react'
import StarterTemplate from '../../components/starter-template/StarterTemplate'
import Topbar from '../../components/topbar/Topbar'
import create from '../../assets/create.svg'
import postgresIcon from '../../assets/postgresIcon.svg'
import mysqlIcon from '../../assets/mysqlIcon.svg'
import mongoIcon from '../../assets/mongoIcon.svg'
import './create-project.css'
import { connect } from 'react-redux'
import { Row, Col, Button, Form, Input, Icon } from 'antd'
import { createProject } from "../../actions/index"

class CreateProject extends Component {
  constructor(props) {
    super(props)
    this.state = {
      selected: "sql-postgres"
    };
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
        <Topbar save="false" />
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
                <StarterTemplate icon={postgresIcon} onClick={() => this.handleSelect("sql-postgres")}
                  heading="POSTGRESQL" desc="Lorem ipsum dolor sit amet, cotetur adipiscing elit. Lorem ipsum."
                  recommended={false} selected={this.state.selected}
                  active={this.state.selected === "sql-postgres"} />
              </Col>

              <Col span={6}>
                <StarterTemplate icon={mongoIcon} onClick={() => this.handleSelect("mongo")}
                  heading="MONGODB" desc="Lorem ipsum dolor sit amet, cotetur adipiscing elit."
                  recommended={false} selected={this.state.selected}
                  active={this.state.selected === "mongo"} />
              </Col>

              <Col span={6}>
                <StarterTemplate icon={mysqlIcon} onClick={() => this.handleSelect("sql-mysql")}
                  heading="MYSQL" desc="Lorem ipsum dolor sit amet, cotetur adipiscing elit. Lorem ipsum dolor."
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
    handleNext: createProject
  }
}

export default connect(mapDispatchToProps)(WrappedCreateProject);
