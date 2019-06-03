import React from 'react';
import { Router, Route } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'


export default () => {
  return (
    <Router history={history}>
      <Route exact path="/" component={Login} />
    </Router>
  )
}