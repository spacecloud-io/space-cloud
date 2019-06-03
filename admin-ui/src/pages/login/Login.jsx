import React from 'react';
import { connect } from 'react-redux';

class Login extends React.Component {
  constructor(props) {
    super(props);
  }

  render() {
    return (
      <div>Login</div>
    )
  }
}

const mapStateToProps = (state) => {
  return {

  }
}

const mapDispatchToProps = (dispatch) => {
  return {

  }
}

export default connect(mapStateToProps, mapDispatchToProps)(Login);