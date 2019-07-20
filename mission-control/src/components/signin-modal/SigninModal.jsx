import React from 'react'
import { connect } from "react-redux";
import { get, set } from 'automate-redux';
import * as firebase from "firebase/app";
import "firebase/auth";
import service from '../../index';
import store from "../../store"
import { notify, handleSpaceUpLoginSuccess } from "../../utils"

import { Modal } from 'antd';
import googleIcon from "../../assets/google.svg"
import githubIcon from "../../assets/github.svg"
import twitterIcon from "../../assets/twitter.svg"
import "./signin-modal.css"

function SigninModal({ handleOauthLogin, visible, handleCancel }) {
  return (
    <Modal
      style={{}}
      wrapClassName="signin-modal"
      closable={false}
      title=""
      footer={null}
      visible={visible}
      onCancel={handleCancel}
    >
      <div>
        <div className="background" />
        <div className="content">
          <h1>Signin to upgrade!</h1>
          <p className="desc">With Space Cloud Enterprise you get the power to experience Firebase + Heroku on your Kubernetes cluster and much more. Signin to access the Space Cloud Enterprise and unlock new powers!</p>
          <div className="footer">
            <img src={googleIcon} alt="" onClick={() => handleOauthLogin("google")} />
            <img src={githubIcon} alt="" onClick={() => handleOauthLogin("twitter")} />
            <img src={twitterIcon} alt="" onClick={() => handleOauthLogin("github")} />
          </div>
        </div>
      </div>
    </Modal>
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
          const newOperationConfig = Object.assign({}, get(store.getState(), "operationConfig", {}), { userId: user.id, key: user.key })
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

