import React from 'react';
import { Form, Input, Select, Switch } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function DeployForm(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="configure">
      <Form className="conn-form" layout="inline">
        <div className="conn-form-options">
          <Form.Item label="Enabled" className="switch">
            {getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
          </Form.Item>
          <Form.Item>
            {getFieldDecorator('orchestrator', {
              rules: [{ required: true, message: '' }]
            })(
              <Select placeholder="Orchestration tool" className="select">
                <Option value="kubernetes">Kubernetes</Option>
              </Select>
            )}
          </Form.Item>
        </div>
        <div className="conn-form-cert">
          <Form.Item className="conn-form-cert-input">
            {getFieldDecorator('namespace', {
              rules: [{ required: true, message: '' }]
            })(<Input style={{ width: 600 }} placeholder="Enter Namespace" />)}
          </Form.Item>
        </div>
        <div className="conn-form-cert">
          <Form.Item className="conn-form-cert-input">
            {getFieldDecorator('registry.url', {
              rules: [{ required: true, message: '' }]
            })(<Input style={{ width: 600 }} placeholder="Enter URL" />)}
          </Form.Item>
        </div>
        <div className="conn-form-cert">
          <Form.Item className="conn-form-cert-input">
            {getFieldDecorator('registry.id', {
              rules: [{ required: true, message: '' }]
            })(<Input style={{ width: 600 }} placeholder="Enter ID" />)}
          </Form.Item>
        </div>
        <div className="conn-form-cert">
          <Form.Item className="conn-form-cert-input">
            {getFieldDecorator('registry.key', {
              rules: [{ required: true, message: '' }]
            })(<Input style={{ width: 600 }} placeholder="Enter Key" />)}
          </Form.Item>
        </div>
      </Form>
    </div>
  );
}

const WrappedDeployForm = Form.create({
  mapPropsToFields(props) {
    return {
      orchestrator: createFormField({ value: props.formState.orchestrator }),
      enabled: createFormField({ value: props.formState.enabled }),
      namespace: createFormField({ value: props.formState.namespace }),
      registry: {
        url: createFormField({ value: props.formState.registry.url }),
        id: createFormField({ value: props.formState.registry.id }),
        key: createFormField({ value: props.formState.registry.key })
      }
    };
  },
  onValuesChange(props, changedValues, allValues) {
    props.handleChange(allValues);
  }
})(DeployForm);

export default WrappedDeployForm;
