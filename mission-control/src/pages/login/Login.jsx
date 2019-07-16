import React from 'react'
import './login.css'
import LoginForm from './LoginForm'
import { Row, Col } from 'antd'
import logo from '../../assets/logo-black.svg'
import loginBg from '../../assets/login.svg'
import { connect } from 'react-redux'
import { set } from "automate-redux"
import { login } from "../../actions/index"

function Login(props) {
  return (
    <div className="login">
      <div className="main-wrapper">
        <Row className="row">
          <Col span={12} className="left-wrapper">
            <div className="left-content">
              <img className="logo" src={logo} alt="logo" /><br />
              <div className="welcome">Welcome back!</div>
              <div className="text">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod.</div><br />
              <img src={loginBg} alt="login" height="240" width="360" /><br />
            </div>
          </Col>

          <Col span={12} className="right-wrapper">
            <div className="right-content">
              <LoginForm formState={props.formState} isLoading={props.isLoading}
                updateFormState={props.updateFormState} handleSubmit={props.handleSubmit} />
            </div>
          </Col>
        </Row>
      </div>
    </div>
  )
}

const mapStateToProps = (state) => {
  return {
    formState: state.uiState.login.formState,
    isLoading: state.uiState.login.isLoading,
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    updateFormState: (fields) => {
      dispatch(set("uiState.login.formState", fields))
    },
    handleSubmit: login
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Login);
