import React from 'react'
import { Modal, Icon, Button, Table } from 'antd'
import Header from "../../components/header/Header"
import { Link } from "react-router-dom"
import './select-project.css'
import { connect } from 'react-redux'
import { deleteProject, openProject } from '../../actions';
import { get } from 'automate-redux';

function SelectProject(props) {
  const columns = [
    {
      title: '',
      dataIndex: 'selected',
      key: 'selected',
      render: (_, record) => {
        return (
          <div>
            {record.selected && <Icon type="check" className="checked" />}
          </div>
        )
      },

      onCell: (record, _) => {
        return {
          selected: record.selected
        };
      }
    },
    { title: 'Project Name', dataIndex: 'name', key: 'projectName' },
    { title: 'ID', dataIndex: 'projectId', key: 'projectId' },
    {
      title: '',
      dataIndex: '',
      key: 'x',
      render: (_, record) => <a onClick={(e) => {
        e.stopPropagation()
        props.handleDelete(record.projectId)
      }
      } href="javascript:;">Delete</a>,
    },
  ];

  const projects = props.projects.map(project => Object.assign({}, project, { selected: project.projectId === props.projectId }))
  return (
    <div >
      <Modal className="select-project" footer={null} closable={false} bodyStyle={{ widtht: "800" }}
        title={<div className="modal-header">
          <Header name="Select a project" />
          <Link to="/mission-control/create-project">
            <Button type="primary" >Create a project</Button>
          </Link>
        </div>}
        visible={props.visible}
        onCancel={props.handleCancel}
      >
        <Table
          pagination={false}
          columns={columns}
          size="middle"
          dataSource={projects}

          onRow={(record) => {
            return {
              onClick: () => {
                {
                  if (!record.selected) {
                    props.handleProjectChange(record.projectId)
                    props.handleCancel()
                  }
                }
              }
            };
          }}
        />
      </Modal>
    </div>
  )
}

const mapStateToProps = (state, ownProps) => {
  return {
    projectId: get(state, "config.id", ""),
    projects: get(state, "projects", []).map(obj => Object.assign({}, { projectId: obj.id, name: obj.name })),
    visible: ownProps.visible
  }
}

const mapDispatchToProps = (dispatch, ownProps) => {
  return {
    handleDelete: deleteProject,
    handleProjectChange: openProject,
    handleCancel: ownProps.handleCancel
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(SelectProject);

