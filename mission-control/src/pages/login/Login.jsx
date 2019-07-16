import React from 'react';
import { connect } from 'react-redux';
import { set } from "automate-redux";
import service from '../../index';
import { notify, handleClusterLoginSuccess } from "../../utils"

import Header from '../../components/header/Header'
import LoginForm from './LoginForm';
import background from '../../assets/Background.svg';
import './login.css';

function Login(props) {
  return (
    <div className="wrapper" >
      <div className="image">
        <img className="background" src={background} alt="background" />
        <div className="r-text">
          <Header name="Welcome to Space Cloud !" color="#FFF" fontSize="28px" />
          <p className="text">Develop great applications without having to write backend code.
              Focus more on business and less on technology.</p>
        </div>
      </div>
      <LoginForm formState={props.formState} isLoading={props.isLoading} updateFormState={props.updateFormState} handleSubmit={props.handleSubmit} />
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
