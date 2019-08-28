import React from 'react';
import './configure.css';
import { Form, Input, Select, Switch } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function FunctionConfigure(props) {
	const { getFieldDecorator } = props.form;
	return (
		<div className="configure">
			<div className="conn-string">Function : </div>

			<Form className="conn-form" layout="inline">
				<div className="conn-form-options">
					<Form.Item>
						{getFieldDecorator('broker', {
							rules: [ { required: true, message: '' } ]
						})(
							<Select placeholder="Broker" className="select">
								<Option value="nats">NATS</Option>
							</Select>
						)}
					</Form.Item>
					<Form.Item label="Enabled" className="switch">
						{getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
					</Form.Item>
				</div>
				<div className="conn-form-cert">
					<Form.Item className="conn-form-cert-input">
						{getFieldDecorator('conn', {
							rules: [ { required: true, message: '' } ]
						})(<Input style={{ width: 600 }} placeholder="Enter Connection String" />)}
					</Form.Item>
				</div>
			</Form>
		</div>
	);
}

const WrappedFunctionConfigureForm = Form.create({
	mapPropsToFields(props) {
		return {
			broker: createFormField({ value: props.formState.broker }),
			enabled: createFormField({ value: props.formState.enabled }),
			conn: createFormField({ value: props.formState.conn })
		};
	},
	onValuesChange(props, changedValues) {
		props.handleChange(changedValues);
	}
})(FunctionConfigure);

export default WrappedFunctionConfigureForm;
