import React from 'react';
import './configure.css';
import { Form, Input, Switch } from 'antd';
import { createFormField } from 'rc-form';

function SslConfigure(props) {
	const { getFieldDecorator } = props.form;
	return (
		<div className="configure">
			<div className="conn-string">SSL : </div>

			<Form className="conn-form" layout="inline">
				<div className="conn-form-switch">
					<Form.Item label="Enabled" className="switch">
						{getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
					</Form.Item>
				</div>
				<div className="conn-form-cert">
					<Form.Item className="conn-form-cert-input">
						{getFieldDecorator('cert', {
							rules: [ { required: true, message: '' } ]
						})(<Input style={{ width: 600 }} placeholder="Enter Certificate" />)}
					</Form.Item>
				</div>
        <div className="conn-form-key">
					<Form.Item className="conn-form-cert-input">
						{getFieldDecorator('key', {
							rules: [ { required: true, message: '' } ]
						})(<Input style={{ width: 600 }} placeholder="Enter Key" />)}
					</Form.Item>
				</div>
			</Form>
		</div>
	);
}

const WrappedSslConfigureForm = Form.create({
	mapPropsToFields(props) {
		return {
			cert: createFormField({ value: props.formState.cert }),
			enabled: createFormField({ value: props.formState.enabled }),
			key: createFormField({ value: props.formState.key })
		};
	},
	onValuesChange(props, changedValues) {
		props.handleChange(changedValues);
	}
})(SslConfigure);

export default WrappedSslConfigureForm;
