import React from 'react'
import logo from '../../assets/logo-black.svg';
import { Form, Icon, Input, Button } from 'antd';
import { createFormField } from 'rc-form';

function LoginForm(props) {
  const handleSubmit = e => {
    e.preventDefault();
    props.form.validateFields((err, values) => {
      if (!err) {
        props.handleSubmit(values.userName, values.password);
      }
    });
  }
  const { getFieldDecorator } = props.form;
  
  return (
    <div className="login-card">
      <div className="content">
        <img className="logo" src={logo} alt="logo"/><br />
        <p className="about">If you've got an idea, let's space it up !</p>

        <Form onSubmit={handleSubmit} className="login-form">
          <Form.Item className="input">
            {getFieldDecorator('userName', {
              rules: [{ required: true, message: 'Please input your username!' }],
            })(
              <Input
                prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />}
                placeholder="Username"

              />,
            )}
          </Form.Item>

          <Form.Item>
            {getFieldDecorator('password', {
              rules: [{ required: true, message: 'Please input your password!' }],
            })(
              (<Input.Password
                prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />}
                type="password"
                placeholder="Password"
              />)
            )}
          </Form.Item>

          <Button type="primary" htmlType="submit" className="login-button" loading={props.isLoading}>
            Log in
          </Button>
        </Form>

      </div>
    </div>
  )
}

const WrappedNormalLoginForm = Form.create({
  mapPropsToFields(props) {
    return {
      userName: createFormField(props.formState.userName),
      password: createFormField(props.formState.password),
    };
  },
  onFieldsChange(props, changedFields, allFields) {
    props.updateFormState(allFields)
  },
})(LoginForm);

export default WrappedNormalLoginForm
