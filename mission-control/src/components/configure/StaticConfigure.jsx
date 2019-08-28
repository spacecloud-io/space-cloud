import React from 'react';
import './configure.css';
import { Form, Input, Select, Switch } from 'antd';
import { createFormField } from 'rc-form';
const { Option } = Select;

function StaticConfigureForm(props) {
	const { getFieldDecorator } = props.form;
	return (
		<div className="configure">
			<div className="conn-string">Static : </div>
			<Form className="conn-form" layout="inline">
				<div className="static">
					<div className="conn-form-options">
						<Form.Item label="Enabled" className="switch static-switch">
							{getFieldDecorator('enabled', { valuePropName: 'checked' })(<Switch size="small" />)}
						</Form.Item>
					</div>
				</div>
			</Form> <br />
		</div>
	);
}

const WrappedStaticConfigureForm = Form.create({
	mapPropsToFields(props) {
		return {
			enabled: createFormField({ value: props.formState.enabled })
		};
	},
	onValuesChange(props, changedValues) {
		props.handleChange(changedValues);
	}
})(StaticConfigureForm);

export default WrappedStaticConfigureForm;
