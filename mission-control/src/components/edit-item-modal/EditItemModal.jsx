import React from 'react'
import { Modal } from 'antd';
import Header from '../header/Header';
import './edit-item-modal.css'
import { Form, Input, Button } from 'antd';

function EditItemModal(props) {
  const handleSubmit = e => {
    e.preventDefault();
    props.form.validateFields((err, values) => {
      if (!err) {
        props.handleSubmit(values.item);
        props.handleCancel();
        props.form.resetFields();
      }
    });
  }
  const { getFieldDecorator } = props.form;

  return (
    <div>
      <Modal className="edit-item-modal" footer={null}
        title={props.heading}
        visible={props.visible}
        onCancel={props.handleCancel}
      >

        <div className="modal-flex">
          <div className="content">
            <Header name={props.name} color="#000" fontSize="18px" />
            <span className="desc">{props.desc}</span>
            <Form onSubmit={handleSubmit} className="edit-form">
              <Form.Item>
                {getFieldDecorator('item', {
                  rules: [{ required: true, message: 'Please input a valid value!'}],
                  initialValue: props.initialValue
                })(
                  <Input className="input"
                    placeholder={props.placeholder}
                  />,
                )}
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit" className="button">
                  DONE
                </Button>
              </Form.Item>
            </Form>
          </div>

          <div className="graphics">
            <img className="vector" src={props.graphics} alt="vector"/><br />
          </div>
        </div>
      </Modal>
    </div>
  );
}

const WrappedNormalLoginForm = Form.create({ name: 'normal_login' })(EditItemModal);

export default WrappedNormalLoginForm

