import React, { useState } from 'react'
import { Form, Input, Button, Tag } from 'antd';
import { Controlled as CodeMirror } from 'react-codemirror2';
import 'codemirror/theme/material.css';
import 'codemirror/lib/codemirror.css';
import 'codemirror/mode/javascript/javascript'
import 'codemirror/addon/selection/active-line.js'
import 'codemirror/addon/edit/matchbrackets.js'
import 'codemirror/addon/edit/closebrackets.js'
import "./testing-form.css"
import service from '../../../index';
import { notify } from "../../../utils"

function TestingForm(props) {
  const [params, setParams] = useState("{}")
  const [loading, setLoading] = useState(null)
  const [respone, setResponse] = useState(null)
  const handleSubmit = e => {
    e.preventDefault();
    props.form.validateFields((err, values) => {
      if (!err) {
        let paramsObj = {}
        try {
          paramsObj = JSON.parse(params)
        } catch (error) {
          notify("error", "Error", "Error in params format")
          return
        }
        setLoading(true)
        service.triggerFunction(props.projectId, values.service, values.function, paramsObj).then(res => {
          setResponse(res)
          setLoading(false)
        })
      }
    });
  };
  const { getFieldDecorator } = props.form;
  return (
    <div className="testing-form">
      <Form onSubmit={handleSubmit} layout="inline">
        <div className="row input-row">
          <Form.Item>
            {getFieldDecorator('service', {
              rules: [{ required: true, message: 'Please select a service!' }],
            })(
              <Input placeholder="Service name" />
            )}
          </Form.Item>
          <Form.Item>
            {getFieldDecorator('function', {
              rules: [{ required: true, message: 'Please select a function!' }],
            })(
              <Input placeholder="Function name" />
            )}
          </Form.Item>
          <Button type="primary" htmlType="submit" loading={props.isLoading} className="btn">
            Trigger function
          </Button>
        </div>
        <div className="params row">
          Function Params: 
          <CodeMirror
            value={params}
            options={{
              mode: { name: "javascript", json: true },
              lineNumbers: true,
              styleActiveLine: true,
              matchBrackets: true,
              autoCloseBrackets: true,
              tabSize: 2,
              autofocus: true
            }}
            onBeforeChange={(editor, data, value) => {
              setParams(value)
            }}
          />
        </div>
        {loading === true && <span>Loading...</span>}
        {loading === false && <div className="row">
          <h3>Status: <Tag color={respone.status === 200 ? "#4F8A10" : "#D8000C"}>{respone.status}</Tag></h3>
          <h3>Result: </h3>
          <div className="result">
            <CodeMirror
              value={JSON.stringify(respone.data.result, null, 2)}
              options={{
                mode: { name: "javascript", json: true },
                tabSize: 2,
                readOnly: true
              }}
              onBeforeChange={(editor, data, value) => {
                setParams(value)
              }}
            />
          </div>
        </div>}
      </Form>
    </div>
  )
}

const WrappedTestingForm = Form.create({})(TestingForm);

export default WrappedTestingForm
