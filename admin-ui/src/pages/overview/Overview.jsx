import React from 'react'
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

function Overview(props) {
  return (
    <div className="overview">
      <Topbar title="Project Overview" />
      <div className="flex-box">
        <Sidenav selectedItem="overview" />
        <div className="page-content ">
          <Row>
            <Col span={9}>
              <Header name="Get started by adding Space Cloud to your app" color="#000" fontSize="22px" />
            </Col>
          </Row>
          <div className="desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. </div>
          <div className="lang">
            <a href="https://www.spaceuptech.com" target="_blank" rel="noopener noreferrer">
              <span class="circle"><img src={python} alt="python" /></span>
            </a>
            <a href="https://www.spaceuptech.com" target="_blank" rel="noopener noreferrer">
              <span class="circle"><img src={js} alt="js" /></span>
            </a>
            <a href="https://www.spaceuptech.com" target="_blank" rel="noopener noreferrer">
              <span class="circle"><img src={java} alt="java" /></span>
            </a>
            <a href="https://www.spaceuptech.com" target="_blank" rel="noopener noreferrer">
              <span class="circle" id="go"><img src={go} alt="go" /></span>
            </a>
            <div className="sepration"></div>
            <a href="https://www.youtube.com" target="_blank" rel="noopener noreferrer">
              <Button type="primary" shape="round" icon="play-circle" size="large" className="get-started">Getting Started</Button>
            </a>
          </div>
          <Header name="Explore Modules" color="#000" fontSize="22px" />
          <Row>
            <Link to={`/${props.projectId}/user-management`}>
              <Col span={11}><UserManagement modules={props.modules.userManagement} /></Col>
            </Link>
            <Link to={`/${props.projectId}/database`}>
              <Col span={11} offset={2}><Database modules={props.modules.database} /></Col>
            </Link>
            <Link to={`/${props.projectId}/functions`}>
              <Col span={11}><Functions modules={props.modules.functions} /></Col>
            </Link>
            <Link to={`/${props.projectId}/configure`}>
              <Col span={11} offset={2}><Configure modules={props.modules.configure} /></Col>
            </Link>
          </Row>
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    projectId: "my-project",
    modules: {
      userManagement: {
        enabled: false,
        mail: true,
        google: true,
        fb: true,
        twitter: true,
        github: true,
      },
      database: {
        enabled: false,
        mysql: true,
        postgres: true,
        mongo: true,
      },
      functions: {
        enabled: false,
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
