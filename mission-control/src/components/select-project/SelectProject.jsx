import React from 'react'
import { connect } from 'react-redux'
import { get, set, reset } from 'automate-redux';
import service from "../../index";
import store from "../../store"
import history from "../../history";
import { openProject, notify, openPlansPage } from "../../utils"

import { Modal, Icon, Button, Table } from 'antd'
import Header from "../../components/header/Header"
import './select-project.css'

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
          <Button type="primary" onClick={props.handleCreateProject}>Create a project</Button>
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
    handleCreateProject: () => {
      const mode = get(store.getState(), "operationConfig.mode", 0)
      if (mode < 1) {
        openPlansPage()
        notify("info", "Info", "You need to upgrade to create multiple projects on the same cluster")
        return
      }
      history.push("/mission-control/create-project")
    },
    handleDelete: (projectId) => {
      service.deleteProject(projectId).then(() => {
        notify("success", "Success", "Project deleted successfully")
        const updatedProjects = get(store.getState(), "projects", []).filter(project => project.id !== projectId)
        dispatch(set("projects", updatedProjects))
        const selectedProject = get(store.getState(), "config.id")
        if (selectedProject === projectId) {
          dispatch(reset("config"))
          dispatch(reset("savedConfig"))
          if (updatedProjects.length) {
            openProject(updatedProjects[0].id)
            return
          }
          history.push("/mission-control/welcome")
        }
      }).catch(ex => {
        console.log("Error", ex)
        notify("error", "Error", ex)
      })
    },
    handleProjectChange: openProject,
    handleCancel: ownProps.handleCancel
  }
}

export default connect(mapStateToProps, mapDispatchToProps)(SelectProject);

