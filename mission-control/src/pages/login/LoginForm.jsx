import React from 'react'
import { Form, Input, Button } from 'antd';
import { createFormField } from 'rc-form';
import './login.css'

function LoginForm(props) {
  const handleSubmit = e => {
    e.preventDefault();
    props.form.validateFields((err, values) => {
      if (!err) {
        props.handleSubmit(values.userName, values.password);
      }
    });
  };
  const { getFieldDecorator } = props.form;
  return (
    <div className="login-form">
      <p className="sign-in">Sign In</p>
      <Form onSubmit={handleSubmit}>
        <Form.Item >
          {getFieldDecorator('userName', {
            rules: [{ required: true, message: 'Please input your username!' }],
          })(
            <Input placeholder="Username" className="input" />,
          )}
        </Form.Item>

        <Form.Item>
          {getFieldDecorator('password', {
            rules: [{ required: true, message: 'Please input your password!' }],
          })(
            <Input.Password type="password" placeholder="Password" className="input" />
          )}
        </Form.Item>

        <Button type="primary" htmlType="submit" loading={props.isLoading} className="btn">
          SIGN IN
          </Button>
      </Form>
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
  onFieldsChange(props, _, allFields) {
    props.updateFormState(allFields)
  },
})(LoginForm);

export default WrappedNormalLoginForm
