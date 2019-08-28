import React from 'react'
import { Form, Switch, Input } from 'antd';
import { createFormField } from 'rc-form';
import { CopyToClipboard } from 'react-copy-to-clipboard';

function Oauth(props) {
  const { getFieldDecorator } = props.form;
  return (
    <div className="oauth">
      <span className="oauth-enable">
        <Form hideRequiredMark={true}>
          <div className="en-button">
            <span>Enable:  </span>
            <Form.Item>
              {getFieldDecorator('enabled', { valuePropName: 'checked' })(
                <Switch size="small" className="en-switch" />
              )}
            </Form.Item>
          </div>

          <div className="details">
            <Form.Item label="App Id">
              {getFieldDecorator('appId', {
                rules: [{ required: true, message: 'Please input your app Id!' }],
              })(
                <Input
                  placeholder="AppId"
                />,
              )}
            </Form.Item>
            <Form.Item label="App Secret">
              {getFieldDecorator('appSecret', {
                rules: [{ required: true, message: 'Please input your app Secret!' }],
              })(
                <Input
                  placeholder="App Secret"
                />,
              )}
            </Form.Item>
            <span>To complete setup, add this OAuth redirect URL to your {props.type} app configuration.</span>

            <div className="url">
              <span>{props.redirectUrl}</span>
              <CopyToClipboard text={props.redirectUrl} >
                <i className="material-icons copy">content_copy</i>
              </CopyToClipboard>
            </div>

          </div> <br />
        </Form>
      </span>
    </div>
  )
}

const WrappedOauthForm = Form.create({
  mapPropsToFields(props) {
    return {
      enabled: createFormField({ value: props.formState.enabled }),
      appId: createFormField({ value: props.formState.appId }),
      appSecret: createFormField({ value: props.formState.appSecret }),
    };
  },
  onValuesChange(props, changedValues) {
    props.handleChange(changedValues)
  },
})(Oauth);

export default WrappedOauthForm


