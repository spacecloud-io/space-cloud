import React, { useState } from 'react'
import ReactGA from "react-ga"
import { connect } from "react-redux";
import { get, set } from 'automate-redux';
import service from '../../index';
import store from "../../store"
import { notify, handleSpaceUpLoginSuccess } from "../../utils"

import { Modal, Icon, Tabs } from 'antd';
import WrappedLoginForm from "./LoginForm"
import WrappedRegisterForm from "./RegisterForm"
import "./signin-modal.css"

const TabPane = Tabs.TabPane;

const SigninModal = ({ handleCancel, handleLogin, handleRegister, visible }) => {
  useState(() => {
    ReactGA.pageview("/signup");
  }, [])
  return (
    <Modal
      title=""
      visible={visible}
      footer={null}
      onCancel={handleCancel}
      className="signin"
    >
      <div className="card-container signin-content">
        <Tabs type="card">
          <TabPane className="login light-blue" tab="Register" key="0">
            <h2 className="greeting">
              Join the Space Up family! <Icon type="heart" theme="filled" className="heart-emoji heart" />
            </h2>
            <WrappedRegisterForm handleRegister={handleRegister} />
          </TabPane>
          <TabPane className="login light-blue" tab="Login" key="1">
            <h2 className="greeting">
              Join the Space Up family! <Icon type="heart" theme="filled" className="heart-emoji heart" />
            </h2>
            <WrappedLoginForm handleLogin={handleLogin} />
          </TabPane>
        </Tabs>
      </div>
    </Modal>
  )
}

const mapStateToProps = (state) => {
  return {
    visible: get(state, "uiState.isSigninModalVisible", false)
  }
}
const mapDispatchToProps = (dispatch) => {
  return {
    handleCancel: () => dispatch(set("uiState.isSigninModalVisible", false)),
    handleLogin: (email, pass) => {
      service.spaceUpLogin(email, pass).then(({ token, user }) => {
        localStorage.setItem("space-up-token", token)
        const newOperationConfig = Object.assign({}, get(store.getState(), "operationConfig", {}), { userId: user.id, key: user.key })
        service.saveOperationConfig(newOperationConfig).then(() => dispatch(set("operationConfig", newOperationConfig)))
        handleSpaceUpLoginSuccess(token)
        notify("success", "Success", "Login successful")
        dispatch(set("uiState.isSigninModalVisible", false))
      }).catch(error => notify("error", "Login failed", error))
    },
    handleRegister: (name, email, pass) => {
      service.spaceUpRegister(name, email, pass).then(({ token, user }) => {
        localStorage.setItem("space-up-token", token)
        const newOperationConfig = Object.assign({}, get(store.getState(), "operationConfig", {}), { userId: user.id, key: user.key })
        service.saveOperationConfig(newOperationConfig).then(() => dispatch(set("operationConfig", newOperationConfig)))
        handleSpaceUpLoginSuccess(token)
        notify("success", "Success", "Signup successful")
        dispatch(set("uiState.isSigninModalVisible", false))
      }).catch(error => notify("error", "Signup failed", error))
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(SigninModal)

