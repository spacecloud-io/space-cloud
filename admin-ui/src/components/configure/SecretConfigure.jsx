import React from 'react'
import './SecretConfigure.css'
import { Form, Input } from 'antd';
import { createFormField } from 'rc-form';

function SecretConfigure(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="secret-configuration">
      <div className="conn-string">App secret</div>
      <Form className="conn-form" layout="inline">
        <Form.Item>
          {getFieldDecorator('conn', {
            rules: [{ required: true, message: 'Please input a secret !' }],
          })(
            <Input.Password style={{ width: 350 }}
              placeholder="Enter App Secret"
            />,
          )}
        </Form.Item>
      </Form>
    </div>
  )
}

const WrappedSecretConfigureForm = Form.create({
  mapPropsToFields(props) {
    return {
      conn: createFormField({ value: props.formState }),
    };
  },
  onValuesChange(props, changedValues) {
    props.handleChange(changedValues)
  },
})(SecretConfigure);

export default WrappedSecretConfigureForm