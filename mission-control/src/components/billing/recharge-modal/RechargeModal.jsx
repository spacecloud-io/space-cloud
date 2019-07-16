import React from 'react'

import { Modal, Button } from 'antd'
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
          visible={this.props.visible}
          onCancel={this.props.handleCancel}
        >
          <RechargeCard active={this.state.selectedOption === 0} amount={10} extraCredits={0} handleClick={() => this.setState({ selectedOption: 0 })} />
          <RechargeCard active={this.state.selectedOption === 1} amount={100} extraCredits={10} handleClick={() => this.setState({ selectedOption: 1 })} />
          <RechargeCard active={this.state.selectedOption === 2} amount={1000} extraCredits={150} handleClick={() => this.setState({ selectedOption: 2 })} />
          <Button onClick={this.handleSubmit}>Next</Button>
        </Modal>
      </div>
    );
  }
}


export default RechargeModal