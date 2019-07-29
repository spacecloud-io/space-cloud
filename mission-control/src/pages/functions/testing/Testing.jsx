import React, { useState } from 'react';
import ReactGA from 'react-ga';
import { connect } from 'react-redux';
import '../../../index.css';
import Sidenav from '../../../components/sidenav/Sidenav';
import Topbar from '../../../components/topbar/Topbar';
import Header from '../../../components/header/Header';
import TestingForm from "./TestingForm"
import Tabs from "../../../components/functions/tabs/Tabs"
import Documentation from '../../../components/documentation/Documentation';
import { get, set } from "automate-redux";
import store from "../../../store"
import "../functions.css"
import service from '../../../index';
import { resolve } from 'path';

const Testing = (props) => {
  useState(() => {
    ReactGA.pageview("/projects/functions/rules");
  }, [])

  return (
    <div className="functions-content">
      <Topbar showProjectSelector />
      <div className="flex-box">
        <Sidenav selectedItem="functions" />
        <div className="page-content">
          {/* <div className="header-flex">
					</div> */}
          <Tabs activeKey="testing" projectId={props.projectId} />
          <div className="documentation-container">
            <Documentation url="https://spaceuptech.com/docs/functions" />
          </div>
          <TestingForm projectId={props.projectId}/>
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state, ownProps) => {
  return {
    rules: get(state, `config.modules.functions.rules`, {}),
    projectId: get(state, "config.id", "")
  }
};

const mapDispatchToProps = (dispatch) => {
  return {}
};

export default connect(mapStateToProps, mapDispatchToProps)(Testing);
