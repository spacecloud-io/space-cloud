import React from 'react'
import { connect } from 'react-redux'
import '../../../index.css'
import Sidenav from '../../../components/sidenav/Sidenav'
import Topbar from '../../../components/topbar/Topbar'
import Header from '../../../components/header/Header'
import Documentation from '../../../components/documentation/Documentation'
import DbConfigure from '../../../components/database-rules/DbConfigure'
import EmptyState from '../../../components/rules/EmptyState'
import rules from '../../../assets/rules.svg'

function Rules(props) {
  return (
    <div>
      <Topbar title="Database" />
      <div className="flex-box">
        <Sidenav selectedItem="database" />
        <div className="page-content">
          <div className="header-flex">
            <Header name="Rules" color="#000" fontSize="22px" />
            <Documentation url="https://spaceuptech.com/docs/database" />
          </div>
          <DbConfigure updateFormState={props.updateFormState} formState={props.formState} />
          <EmptyState graphics={rules} desc="Guard your data with rules that define who has access to it and how it is structured." buttonText="Add a table" handleClick={props.handleClick} />
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    formState: {
      enabled: true,
      conn: "http://localhost/8080"
    },
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleClick: () => {
      console.log("Checked")
    },

    updateFormState: (fields) => {
      console.log(fields)
    },
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
