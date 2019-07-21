import React from 'react'
import { Button, Input, Form } from 'antd';
const LoginForm = (props) => {
  const { form, handleLogin } = props
  const { getFieldDecorator } = form;
  const handleSubmit = e => {
    e.preventDefault();
    form.validateFields((err, values) => {
      if (!err) {
        handleLogin(values.email, values.pass);
      }
    });
  }
  return (
    <Form onSubmit={handleSubmit}>
      <Form.Item>
        {getFieldDecorator('email', {
          rules: [{ required: true, type: 'email', message: 'Valid email is required!' }],
        })(
          <Input className="form-input"
            placeholder="Email"
          />,
        )}
      </Form.Item>
      <Form.Item>
        {getFieldDecorator('pass', {
          rules: [{ required: true, message: 'Password is required!' }],
        })(
          <Input className="form-input"
            type="password"
            placeholder="Password"
          />,
        )}
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" className="submit-btn">
          Login
        </Button>
      </Form.Item>
    </Form>
  )
}

const WrappedLoginForm = Form.create({ name: 'login_form' })(LoginForm);

export default WrappedLoginForm

