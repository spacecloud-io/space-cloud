import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from "react-redux";
import store from "./store";
import { handleClusterLoginSuccess, handleSpaceUpLoginSuccess } from './utils';
import Service from "./services/service";
import * as firebase from "firebase/app";

import App from './App';
import * as serviceWorker from './serviceWorker';
import './index.css'




var firebaseConfig = {
  apiKey: "AIzaSyDDk3Nx9Zgft5wfT9oQxJSiObIOYSuIV34",
  authDomain: "space-cloud.firebaseapp.com",
  databaseURL: "https://space-cloud.firebaseio.com",
  projectId: "space-cloud",
  storageBucket: "",
  messagingSenderId: "332138526349",
  appId: "1:332138526349:web:a3c24f2fe681c03e"
};

// Initialize Firebase
firebase.initializeApp(firebaseConfig);

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