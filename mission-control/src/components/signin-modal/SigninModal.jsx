import React from 'react'
import { connect } from "react-redux";
import { get, set } from 'automate-redux';
import * as firebase from "firebase/app";
import "firebase/auth";
import service from '../../index';
import store from "../../store"
import { notify, handleSpaceUpLoginSuccess } from "../../utils"
import { google } from '../../assets/login-google.svg'
import { fb } from '../../assets/login-fb.svg'
import { github } from '../../assets/login-github.svg'
import { twitter } from '../../assets/login-twitter.svg'
import './sign-in-modal.css'

import { Modal, Button, notification } from 'antd';

function SigninModal({ handleOauthLogin, visible, handleCancel }) {
  return (
    <div className="sign-in-modal">
      <div className="orange-bg"></div>
      <div className="blue-bg">
        <span className="modal-heading">Sign up to upgrade!</span><br />
        <span className="modal-desc">Develop great applications without having to write backend code. Focus more
          on business and less on the technology.</span><br />
        <div className="modal-footer">
          <img src={fb} alt="" onClick={() => handleOauthLogin("fb")} />
          <img src={google} alt="" onClick={() => handleOauthLogin("google")} />
          <img src={github} alt="" onClick={() => handleOauthLogin("github")} />
          <img src={twitter} alt="" onClick={() => handleOauthLogin("twitter")} />
        </div>
        {/* <Modal footer={null}
          visible={visible}
          onCancel={handleCancel}
        >
          <img src={fb} alt="" onClick={() => handleOauthLogin("fb")}/>
          <img src={google} alt="" onClick={() => handleOauthLogin("google")}/>
          <img src={github} alt="" onClick={() => handleOauthLogin("github")}/>
          <img src={twitter} alt="" onClick={() => handleOauthLogin("twitter")}/>
        </Modal> */}
      </div>

    </div>
  );
}

const mapStateToProps = (state) => {
  return {
    visible: get(state, "uiState.isSigninModalVisible", false)
  }
}
const mapDispatchToProps = (dispatch) => {
  return {
    handleCancel: () => dispatch(set("uiState.isSigninModalVisible", false)),
    handleOauthLogin: (method) => {
      var provider;
      switch (method) {
        case 'google':
          provider = new firebase.auth.GoogleAuthProvider();
          break;
        case 'fb':
          provider = new firebase.auth.FacebookAuthProvider();
          break;
        case 'twitter':
          provider = new firebase.auth.TwitterAuthProvider();
          break;
        case 'github':
          provider = new firebase.auth.GithubAuthProvider();
          break;
      }
      firebase.auth().signInWithPopup(provider).then(function (result) {
        // The signed-in user info.
        var user = result.user;
        service.oauthLogin(user.uid).then(({ token, user }) => {
          localStorage.setItem("space-up-token", token)
          const newOperationConfig = Object.assign({}, get(store.getState(), "operationConfig", {}), { email: user.email, key: user.key })
          service.saveOperationConfig(newOperationConfig).then(() => dispatch(set("operationConfig", newOperationConfig)))
          handleSpaceUpLoginSuccess(token)
        }).catch(error => {
          console.log("Error", error)
          notify("error", "Error", "Could not signin")
        })
      }).catch(error => {
        console.log("Error", error)
        notify("error", "Error", "Could not signin")
      }).finally(() => dispatch(set("uiState.isSigninModalVisible", false)))
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(SigninModal)

