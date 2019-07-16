import React from 'react'
import './recharge-modal.css'
import { Modal, Button, Row, Col } from 'antd'
import RechargeCard from "../recharge-card/RechargeCard"

class RechargeModal extends React.Component {
  constructor(props) {
    super(props)
    this.state = { selectedOption: 1 }
  }

  handleSubmit = () => {
    this.props.handleCancel()
    this.props.handleSubmit(this.state.selectedOption)
  }

  render() {
    return (
      <div>
        <Modal footer={null}
          className="recharge-modal"
          visible={this.props.visible}
          onCancel={this.props.handleCancel}
          closable={false}
        >
          <span class="pricing-title">Pricing</span>
          <p class="pricing-desc">Lorem ipsum dolor sit amet, consectetur adipiscing elit,
          sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim
          veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
          consequat.</p>
          <Row>
            <Col span={8}>
              <RechargeCard active={this.state.selectedOption === 0} amount={10} extraCredits={0} handleClick={() => this.setState({ selectedOption: 0 })}
                desc="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." />
            </Col>
            <Col span={8}>
              <RechargeCard active={this.state.selectedOption === 1} amount={100} extraCredits={10} handleClick={() => this.setState({ selectedOption: 1 })}
                desc="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." />
            </Col>
            <Col span={8}>
              <RechargeCard active={this.state.selectedOption === 2} amount={1000} extraCredits={150} handleClick={() => this.setState({ selectedOption: 2 })}
                desc="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." />
            </Col>
          </Row>
          <Button type="primary" onClick={this.handleSubmit}>Next</Button>
        </Modal>
      </div>
    );
  }
}


export default RechargeModal