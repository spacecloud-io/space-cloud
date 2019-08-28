import React from 'react';
import { Form, Input, Select, Switch, Divider, Tooltip } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function DeployForm(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="configure">
      <Form className="conn-form" layout="inline">
        <Form.Item label="Enabled" className="switch">
          {getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
        </Form.Item>
        <Divider />
        <div className="deploy-form">
          <div className="deploy-flex">
            <div className="orchestrator">Orchestrator:</div>
            <div className="conn-form-options">
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
          </div>
          <Divider />
          <div className="deploy-flex">
            <div className="namespace">Namespace/Network:</div>
            <div className="conn-form-cert">
              <Form.Item className="conn-form-cert-input">
                {getFieldDecorator('namespace', {
                  rules: [{ required: true, message: '' }]
                })(<Input style={{ width: 400 }} placeholder="Enter Namespace" />)}
              </Form.Item>
            </div>
            <Tooltip placement="right" title="Namespace in kubernetes">
              <span style={{ height: 20 }}><i class="material-icons help">help_outline</i></span>
            </Tooltip>
          </div>
          <Divider />
          <div className="deploy-flex registry">
            <div>Registry:</div>
            <div className="url">URL:</div>
            <div className="conn-form-cert">
              <Form.Item className="conn-form-cert-input">
                {getFieldDecorator('registry.url', {
                  rules: [{ required: true, message: '' }]
                })(<Input style={{ width: 400 }} placeholder="Enter URL" />)}
              </Form.Item>
            </div>
            <Tooltip placement="right" title="URL of the registry">
              <span style={{ height: 20 }}><i class="material-icons help">help_outline</i></span>
            </Tooltip>
          </div>
          <div className="deploy-flex registry">
            <div className="id">ID:</div>
            <div className="conn-form-cert">
              <Form.Item className="conn-form-cert-input">
                {getFieldDecorator('registry.id', {
                  rules: [{ required: true, message: '' }]
                })(<Input style={{ width: 400 }} placeholder="Enter ID" />)}
              </Form.Item>
            </div>
            <Tooltip placement="right" title="User ID">
              <span style={{ height: 20 }}><i class="material-icons help">help_outline</i></span>
            </Tooltip>
          </div>
          <div className="deploy-flex registry">
            <div className="key">KEY:</div>
            <div className="conn-form-cert">
              <Form.Item className="conn-form-cert-input">
                {getFieldDecorator('registry.key', {
                  rules: [{ required: true, message: '' }]
                })(<Input style={{ width: 400 }} placeholder="Enter Key" />)}
              </Form.Item>
            </div>
            <Tooltip placement="right" title="API key">
              <span style={{ height: 20 }}><i class="material-icons help">help_outline</i></span>
            </Tooltip>
          </div>
          <Divider />
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
