import React from 'react'
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
            <span class="circle"><img src={python} alt="python" /></span>
            <span class="circle"><img src={js} alt="js" /></span>
            <span class="circle"><img src={java} alt="java" /></span>
            <span class="circle" id="go"><img src={go} alt="go" /></span>
            <div className="sepration"></div>
            <Button type="primary" shape="round" icon="play-circle" size="large" className="get-started">Getting Started</Button>
          </div>
          <Header name="Explore Modules" color="#000" fontSize="22px" />
          <Row gutter={48}>
            <Col span={11}><UserManagement modules={props.modules.userManagement} /></Col>
            <Col span={11}><Database modules={props.modules.database} /></Col>
            <Col span={11}><Functions modules={props.modules.funcions} /></Col>
            <Col span={11}><Configure modules={props.modules.configure} /></Col>
          </Row>
          
        </div>
      </div>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    modules: {
      userManagement: {
        enabled: true,
        mail: true,
        gmail: true,
        facebook: true,
        twitter: true,
        github: true,
      },
      databse: {
        enabled: true,
        mysql: true,
        postgres: true,
        mongo: true,
      },
      functions: {
        enabled: true,
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
