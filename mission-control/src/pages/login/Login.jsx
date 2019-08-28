import React, { useState } from 'react'
import './login.css'
import { Row, Col } from 'antd'
import logo from '../../assets/logo-black.svg'
import { connect } from 'react-redux'
import { set } from "automate-redux"
import loginBg from '../../assets/login.svg'
import service from '../../index';
import { notify, handleClusterLoginSuccess } from "../../utils"
import LoginForm from './LoginForm';
import ReactGA from 'react-ga';

function Login(props) {
  useState(() => {
    ReactGA.pageview("/");
  }, [])
  return (
    <div className="login">
      <div className="main-wrapper">
        <Row className="row">
          <Col span={12} className="left-wrapper">
            <div className="left-content">
              <img className="logo" src={logo} alt="logo" /><br />
              <div className="welcome">Welcome back!</div>
              <div className="text">Login to configure your Space Cloud cluster.</div><br />
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
    handleSubmit: (user, pass) => {
      service.login(user, pass).then(token => {
        localStorage.setItem("token", token)
        handleClusterLoginSuccess(token)
      }).catch(error => {
        console.log("Error", error)
        notify("error", "Error", "Could not login")
      })
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Login);
