import React from 'react';
import { Router, Route } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'
import Database from './pages/database/Database'
import UserManagement from './pages/user-management/UserManagement';
import Rules from './pages/database/rules/Rules'
import Configure from './pages/configure/configure'
import Overview from  './pages/overview/Overview'

export default () => {
  return (
    <Router history={history}>
      <Route exact path="/" component={Login} />
      <Route exact path="/:projectId/database" component={Database} />
      <Route exact path="/:projectId/user-management" component={UserManagement} />
      <Route exact path="/:projectId/overview" component={Overview} />
      <Route exact path="/:projectId/database/rules/:database" component={Rules} />
      <Route exact path="/:projectId/configure" component={Configure} />

    </Router>
  )
}