import React, { useState } from 'react'
import ReactGA from 'react-ga';
import { Link } from 'react-router-dom'
import { connect } from 'react-redux'
import '../../index.css'
import './overview.css'
import '../../components/card-modules/cards.css'
import Sidenav from '../../components/sidenav/Sidenav'
import Topbar from '../../components/topbar/Topbar'
import Header from '../../components/header/Header'
import python from '../../assets/python.svg'
import js from '../../assets/js.svg'
import java from '../../assets/java.svg'
import go from '../../assets/go.svg'
import UserManagement from '../../components/card-modules/UserManagement'
import Database from '../../components/card-modules/Database'
import Functions from '../../components/card-modules/Functions'
import Configure from '../../components/card-modules/Configure'
import { Row, Col, Button } from 'antd'
import { get } from 'automate-redux';

function Overview(props) {
  useState(() => {
    ReactGA.pageview("/projects/overview");
  }, [])
  return (
    <div className="overview">
      <Topbar showProjectSelector />
      <div className="flex-box">
        <Sidenav selectedItem="overview" />
        <div className="page-content ">
          <Row>
            <Col span={9}>
              <Header name="Get started by adding Space Cloud to your app" color="#000" fontSize="22px" />
            </Col>
          </Row>
          <div className="overview-padding">
            <div className="desc">Start with one of the language specific getting started guides below or watch a getting started video to help you get started.</div>
            <div className="lang">
              <a href="https://www.spaceuptech.com/docs/getting-started/python" target="_blank" rel="noopener noreferrer">
                <span class="circle"><img src={python} alt="python" heg /></span>
              </a>
              <a href="https://www.spaceuptech.com/docs/getting-started/javascript" target="_blank" rel="noopener noreferrer">
                <span class="circle"><img src={js} alt="js" /></span>
              </a>
              <a href="https://www.spaceuptech.com/docs/getting-started/java" target="_blank" rel="noopener noreferrer">
                <span class="circle"><img src={java} alt="java" /></span>
              </a>
              <a href="https://www.spaceuptech.com/docs/getting-started/golang" target="_blank" rel="noopener noreferrer">
                <span class="circle" id="go"><img src={go} alt="go" /></span>
              </a>
              <div className="sepration"></div>
              <a href="https://spaceuptech.com/docs/deployment/" target="_blank" rel="noopener noreferrer">
                <Button type="primary" shape="round" icon="play-circle" size="large" className="get-started">Getting Started</Button>
              </a>
            </div>
          </div>
          <Header name="Explore Modules" color="#000" fontSize="22px" />
          <Row className="overview-padding">
            <Link to={`/mission-control/projects/${props.projectId}/user-management`}>
              <Col span={11}><UserManagement modules={props.modules.userManagement} /></Col>
            </Link>
            <Link to={`/mission-control/projects/${props.projectId}/database`}>
              <Col span={11} offset={1}><Database modules={props.modules.database} /></Col>
            </Link>
            <Link to={`/mission-control/projects/${props.projectId}/functions`}>
              <Col span={11}><Functions modules={props.modules.functions} /></Col>
            </Link>
            <Link to={`/mission-control/projects/${props.projectId}/configure`}>
              <Col span={11} offset={1}><Configure modules={props.modules.configure} /></Col>
            </Link>
          </Row>
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    projectId: ownProps.match.params.projectId,
    modules: {
      userManagement: {
        enabled: get(state, "config.modules.auth.enabled", false),
        email: get(state, "config.modules.auth.email.enabled", false),
        google: get(state, "config.modules.auth.google.enabled", false),
        fb: get(state, "config.modules.auth.fb.enabled", false),
        twitter: get(state, "config.modules.twitter.email.enabled", false),
        github: get(state, "config.modules.auth.github.enabled", false),
      },
      database: {
        enabled: true,
        mysql: get(state, `config.modules.crud.sql-mysql.enabled`, false),
        postgres: get(state, `config.modules.crud.sql-postgres.enabled`, false),
        mongo: get(state, `config.modules.crud.mongo.enabled`, false),
      },
      functions: {
        enabled: get(state, `config.modules.functions.enabled`, false),
      },
      configure: {
        enabled: true,
      },
    }
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Overview);
