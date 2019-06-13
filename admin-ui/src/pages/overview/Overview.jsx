import React from 'react'
import { connect } from 'react-redux'
import '../../index.css'
import Sidenav from '../../components/sidenav/Sidenav'
import Topbar from '../../components/topbar/Topbar'
import Header from '../../components/header/Header'
import Documentation from '../../components/documentation/Documentation'

function Overview(props) {
  return (
      <div className="overview">
      <Topbar title="Project Overview" />
      <div className="flex-box">
        <Sidenav selectedItem="overview" />
        <div className="page-content">
          <div className="header-flex">
            <Header name="Get started by adding Space Cloud to your app" color="#000" fontSize="22px" />
            <Documentation url="https://spaceuptech.com/docs/database" />
          </div>
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Overview);
