import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from "react-redux";
import store from "./store";
import { handleClusterLoginSuccess, handleSpaceUpLoginSuccess } from './utils';
import Service from "./services/service";

import App from './App';
import * as serviceWorker from './serviceWorker';
import './index.css'

import ReactGA from 'react-ga';
if (process.env.NODE_ENV === "production")  {
  ReactGA.initialize('UA-104521884-3');
}


const service = new Service()

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById("root")
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();

export default service

const token = localStorage.getItem("token")
const spaceUpToken = localStorage.getItem("space-up-token")
if (token) {
  const urlParams = window.location.pathname.split("/")
  let lastProjectId = null
  if (urlParams.length > 3 && urlParams[3]) {
    lastProjectId = urlParams[3]
  }
  handleClusterLoginSuccess(token, lastProjectId)
}

if (spaceUpToken) {
  handleSpaceUpLoginSuccess(spaceUpToken)
}