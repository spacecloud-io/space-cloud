import React from 'react';
import { connect } from 'react-redux';
import background from '../../assets/Background.svg';
import Header from '../../components/header/Header'
import './login.css';
import LoginForm from './LoginForm';
import { set } from "automate-redux";

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
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Login);
