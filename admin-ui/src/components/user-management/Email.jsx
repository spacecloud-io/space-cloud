import React from 'react'
import { Form, Switch } from 'antd';
import { createFormField } from 'rc-form';

function Email(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="email">
      <span>No fields to configure</span>
      <span className="email-collapse">
        <span>Enable:  </span>
        <Form>
          <Form.Item>
            {getFieldDecorator('enabled', { valuePropName: 'checked' })(
              <Switch size="small" className="switch" />
            )}
          </Form.Item>
        </Form>
      </span>
    </div>
  )
}

const WrappedEmailForm = Form.create({
  mapPropsToFields(props) {
    return {
      enabled: createFormField({ value: props.formState.enabled }),
    };
  },
  onValuesChange(props, changedValues) {
    props.handleChange(changedValues)
  },
})(Email);

export default WrappedEmailForm
