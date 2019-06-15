import React from 'react';
import { Router, Route } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'
import Database from './pages/database/Database'
import UserManagement from './pages/user-management/UserManagement';

export default () => {
  return (
    <Router history={history}>
      <Route exact path="/" component={Login} />
      <Route exact path="/:projectId/database" component={Database} />
      <Route exact path="/:projectId/user-management" component={UserManagement} />
    </Router>
  )
}