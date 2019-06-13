import React, { Component } from 'react'
// import './rules.css'
import { connect } from 'react-redux'
import { Row ,Col} from 'antd';

class Rules extends Component {
  constructor(props) {
    super(props)
    this.state = {
    };

  }


  render() {
    return (
      <div className="rules-main-wrapper">
      <Row>
        <Col span={6} >
       {this.props.rules.keys.map((rule)=>{
         <p>{rule}</p>
       })}
       </Col>
        <Col span={18} >
        vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv
        </Col>
        </Row>
      </div>
    )
  }
}

const mapStateToProps = (state, ownProps) => {
  return {
  }
}

const mapDispatchToProps = (dispatch) => {
  return {
   
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Rules);
