import React from 'react';
import { connect } from 'react-redux';
import '../../index.css'
import Sidenav from '../../components/sidenav/Sidenav';
import Topbar from '../../components/topbar/Topbar';
import Header from '../../components/header/Header';
import DeployForm from "./DeployForm";
import { get, set } from 'automate-redux';
import EmptyState from "../../components/empty-state/EmptyState"
import store from ".././../store";
import someGraphics from '../../assets/projectId.svg'
import { triggerSignin, isUserSignedIn, openPlansPage } from "../../utils"

function Deploy({ isSignedIn, mode, deployConfig, handleChange }) {
  return (
    <div>
      <Topbar showProjectSelector />
      <div className="flex-box">
        <Sidenav selectedItem="deploy" />
        <div className="page-content">
          {
            isSignedIn ?
              <React.Fragment>
                {mode > 0 ?
                  <React.Fragment>
                    <Header name="Deploy Configuration" color="#000" fontSize="22px" />
                    <DeployForm formState={deployConfig} handleChange={handleChange} />
                  </React.Fragment>
                  :
                  <EmptyState
                    graphics={someGraphics}
                    handleClick={openPlansPage}
                    desc="Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae qui id nulla ipsa maiores fugit ipsum inventore esse iste magnam. Porro blanditiis possimus animi voluptatum? Similique vel illo at asperiores."
                    actionText="Explore Plans" />
                }

              </React.Fragment>
              :
              <EmptyState
                graphics={someGraphics}
                handleClick={triggerSignin}
                desc="Lorem ipsum dolor sit amet consectetur adipisicing elit. Vitae qui id nulla ipsa maiores fugit ipsum inventore esse iste magnam. Porro blanditiis possimus animi voluptatum? Similique vel illo at asperiores."
                actionText="Register / Login" />

          }
        </div>
      </div>
    </div>
  );
}

const mapStateToProps = (state) => {
  return {
    isSignedIn: isUserSignedIn(),
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
