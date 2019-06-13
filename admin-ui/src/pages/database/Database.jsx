import React from 'react'
import { connect } from 'react-redux'
import './database.css'
import '../../index.css'
import Header from '../../components/header/Header'
import mysql from '../../assets/mysql.svg'
import postgresql from '../../assets/postgresql.svg'
import mongodb from '../../assets/mongodb.svg'
import Sidenav from '../../components/sidenav/Sidenav'
import Topbar from '../../components/topbar/Topbar'
import Documentation from '../../components/documentation/Documentation'
import DatabaseCardList from '../../components/database-card/DatabaseCardList'
import { Redirect} from "react-router-dom";

function Database(props) {
  const cards = [{ graphics: mysql, name: "MySQL", desc: "The world's most popular open source database.", key: "sql-mysql"},
  { graphics: postgresql, name: "PostgreSQL", desc: "The world's most popular open source database.", key: "sql-postgres"},
  { graphics: mongodb, name: "MongoDB", desc: "A open-source cross-platform document- oriented database.", key: "mongo"}]

  if(props.selectedDb){
    return <Redirect to={`/${props.projectId}/database/rules/${props.selectedDb}`}/>;
  }
    return (
    <div>
      <Topbar title="Database" />
      <div className="flex-box">
        <Sidenav selectedItem="database" />
        <div className="page-content">
          <div className="header-flex">
            <Header name="Add a database" color="#000" fontSize="22px" />
            <Documentation url="https://spaceuptech.com/docs/database" />
          </div>
          <p className="db-desc">Start using crud by enabling one of the following databases.</p>
          <DatabaseCardList cards={cards} handleEnable={props.handleEnable}/>
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  console.log("OwnProps", ownProps)
  return {
    projectId: "ToDo-App",
    selectedDb: undefined,
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
    handleEnable(key) {
      console.log(`Enabled: ${key}`)
    }
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Database);
