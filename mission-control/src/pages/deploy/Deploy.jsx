import React from 'react';
import { connect } from 'react-redux';
import '../../index.css'
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import Header from '../../components/header/Header';
import Documentation from '../../components/documentation/Documentation'
import DeployForm from "./DeployForm";
import { get, set } from 'automate-redux';
import EmptyState from "../../components/empty-state/EmptyState"
import store from ".././../store";
import someGraphics from '../../assets/projectId.svg'
import { openPlansPage } from "../../utils"
import './deploy.css'

function Deploy({ mode, deployConfig, handleChange }) {
  return (
    <div class="deploy">
      <Topbar showProjectSelector />
      <div className="flex-box">
        <Sidenav selectedItem="deploy" />
        <div className="page-content">
          {mode > 0 ?
            <React.Fragment>
              <div className="header-flex">
                <Header name="Deploy Configuration" color="#000" fontSize="22px" />
                <Documentation url="https://spaceuptech.com" />
              </div>
              <DeployForm formState={deployConfig} handleChange={handleChange} />
            </React.Fragment>
            :
            <EmptyState
              graphics={someGraphics}
              handleClick={openPlansPage}
              desc="Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae qui id nulla ipsa maiores fugit ipsum inventore esse iste magnam. Porro blanditiis possimus animi voluptatum? Similique vel illo at asperiores."
              actionText="Explore Plans" />
          }
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state) => {
  return {
    mode: get(state, "operationConfig.mode", 0),
    deployConfig: get(state, "deployConfig")
  };
};

const mapDispatchToProps = (dispatch) => {
  return {
    handleChange: (values) => {
      dispatch(set("deployConfig", values))
    }
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(Deploy);
