import React from 'react';
import { Router, Route } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'
import Database from './pages/database/Database'
import Rules from './pages/database/rules/Rules'


export default () => {
  return (
    <Router history={history}>
      <Route exact path="/" component={Login} />
      <Route exact path="/:projectId/database" component={Database} />
      <Route exact path="/:projectId/database/rules" component={Rules} />
    </Router>
  )
}