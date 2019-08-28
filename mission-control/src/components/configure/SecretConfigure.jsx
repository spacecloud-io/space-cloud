import React from 'react'
import './configure.css'
import { Form, Input } from 'antd';
import { createFormField } from 'rc-form';

function SecretConfigure(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="configure">
      <div className="conn-string">App secret : </div>
      <Form className="conn-form" layout="inline">
        <Form.Item>
          {getFieldDecorator('secret', {
            rules: [{ required: true, message: 'Please input a secret !' }],
          })(
            <Input.Password style={{ width: 600 }}
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
      secret: createFormField({ value: props.formState }),
    };
  },
  onValuesChange(props, changedValues) {
    props.handleChange(changedValues.secret)
  },
})(SecretConfigure);

export default WrappedSecretConfigureForm