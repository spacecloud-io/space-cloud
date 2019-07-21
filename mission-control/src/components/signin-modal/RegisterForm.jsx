import React from 'react'
import { Button, Input, Form } from 'antd';
const RegisterForm = (props) => {
  const { form, handleRegister } = props
  const { getFieldDecorator } = form;
  const handleSubmit = e => {
    e.preventDefault();
    form.validateFields((err, values) => {
      if (!err) {
        handleRegister(values.name, values.email, values.pass);
      }
    });
  }
  return (
    <Form onSubmit={handleSubmit}>
      <Form.Item>
        {getFieldDecorator('name', {
          rules: [{ required: true, message: 'Name is required!' }],
        })(
          <Input className="form-input"
            placeholder="Name"
          />,
        )}
      </Form.Item>
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
          Create your free account!
                </Button>
      </Form.Item>
    </Form>
  )
}

const WrappedRegisterForm = Form.create({ name: 'register_form' })(RegisterForm);

export default WrappedRegisterForm

