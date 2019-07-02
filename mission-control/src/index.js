import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import * as serviceWorker from './serviceWorker';

import './index.css'

import { Provider } from "react-redux";
import store from "./store";
import Client from "./services/client"
import Service from "./services/service"
import { loadConfig } from "./actions/index";

const client = new Client()
const service = new Service(client)
const token = localStorage.getItem("token")

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

if (token) {
  client.setToken(token)
  const urlParams  = window.location.pathname.split("/")
  if (urlParams.length > 2 && urlParams[2]) {
    loadConfig(urlParams[2])
  }
}