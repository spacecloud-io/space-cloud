import React from 'react';
import { Router, Route, Redirect } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'
import Database from './pages/database/Database'
import UserManagement from './pages/user-management/UserManagement';
import DBRules from './pages/database/rules/Rules'
import FunctionRules from './pages/functions/rules/Rules'
import FileStorageRules from './pages/file-storage/rules/Rules'
import Configure from './pages/configure/configure'
import Overview from  './pages/overview/Overview'

export default () => {
  return (
    <Router history={history}>
      <Route exact path="/mission-control" component={Login} />
      <Route exact path="/mission-control/:projectId" component={(props) => <Redirect to={`/mission-control/${props.match.params.projectId}/overview`}/>} />
      <Route exact path="/mission-control/:projectId/overview" component={Overview} />
      <Route exact path="/mission-control/:projectId/user-management" component={UserManagement} />
      <Route exact path="/mission-control/:projectId/database" component={Database} />
      <Route exact path="/mission-control/:projectId/database/rules/:database" component={DBRules} />
      <Route exact path="/mission-control/:projectId/functions" component={(props) => <Redirect to={`/mission-control/${props.match.params.projectId}/functions/rules`} />} />
      <Route exact path="/mission-control/:projectId/functions/rules" component={FunctionRules} />
      <Route exact path="/mission-control/:projectId/file-storage" component={(props) => <Redirect to={`/mission-control/${props.match.params.projectId}/file-storage/rules`} />} />
      <Route exact path="/mission-control/:projectId/file-storage/rules" component={FileStorageRules} />
      <Route exact path="/mission-control/:projectId/configure" component={Configure} />

    </Router>
  )
}