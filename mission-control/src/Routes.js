import React from 'react';
import { Router, Route, Redirect } from 'react-router-dom';
import history from './history';
import Login from './pages/login/Login'
import Home from "./pages/home/Home"
import Database from './pages/database/Database'
import UserManagement from './pages/user-management/UserManagement';
import DBRules from './pages/database/rules/Rules'
import FunctionRules from './pages/functions/rules/Rules'
import FileStorageRules from './pages/file-storage/rules/Rules'
import StaticRules from './pages/static/rules/Rules'
import Configure from './pages/configure/configure'
import Overview from './pages/overview/Overview'
import Welcome from "./pages/welcome/Welcome"
import CreateProject from "./pages/create-project/CreateProject"
import Billing from "./pages/billing/Billing"
import Deploy from "./pages/deploy/Deploy"
import Plans from "./pages/plans/Plans"
import SigninModal from "./components/signin-modal/SigninModal"
import Explorer from './pages/explorer/Explorer';

export default () => {
  return (
    <Router history={history}>
      <Route exact path="/mission-control" component={Home} />
      <Route exact path="/mission-control/login" component={Login} />
      <Route exact path="/mission-control/welcome" component={Welcome} />
      <Route exact path="/mission-control/create-project" component={CreateProject} />
      <Route exact path="/mission-control/projects/:projectId" component={(props) => <Redirect to={`/mission-control/projects/${props.match.params.projectId}/overview`} />} />
      <Route exact path="/mission-control/projects/:projectId/overview" component={Overview} />
      <Route exact path="/mission-control/projects/:projectId/user-management" component={UserManagement} />
      <Route exact path="/mission-control/projects/:projectId/database" component={Database} />
      <Route exact path="/mission-control/projects/:projectId/database/rules/:database" component={DBRules} />
      <Route exact path="/mission-control/projects/:projectId/functions" component={(props) => <Redirect to={`/mission-control/projects/${props.match.params.projectId}/functions/rules`} />} />
      <Route exact path="/mission-control/projects/:projectId/functions/rules" component={FunctionRules} />
      <Route exact path="/mission-control/projects/:projectId/file-storage" component={(props) => <Redirect to={`/mission-control/projects/${props.match.params.projectId}/file-storage/rules`} />} />
      <Route exact path="/mission-control/projects/:projectId/file-storage/rules" component={FileStorageRules} />
      <Route exact path="/mission-control/projects/:projectId/gateway" component={(props) => <Redirect to={`/mission-control/projects/${props.match.params.projectId}/gateway/rules`} />} />
      <Route exact path="/mission-control/projects/:projectId/gateway/rules" component={StaticRules} />
      <Route exact path="/mission-control/projects/:projectId/configure" component={Configure} />
      <Route exact path="/mission-control/projects/:projectId/explorer" component={Explorer} />
      <Route exact path="/mission-control/projects/:projectId/deploy" component={Deploy} />
      <Route exact path="/mission-control/projects/:projectId/plans" component={Plans} />
      <Route exact path="/mission-control/projects/:projectId/billing" component={Billing} />
      <SigninModal />
    </Router>
  )
}