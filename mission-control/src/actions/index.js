import service from "../index";
import store from "../store";
import history from "../history";
import { set, get } from "automate-redux";
import { cloneDeep } from "lodash"
import { notification } from "antd";
import { adjustConfig, unAdjustConfig } from "./helpers"

export const loadConfig = (project) => {
  service.loadConfig(project).then(({ config }) => {
    const adjustedConfig = adjustConfig(config)
    store.dispatch(set("config", adjustedConfig))
    store.dispatch(set("savedConfig", cloneDeep(adjustedConfig)))
  })
}

export const login = (user, pass) => {
  service.login(user, pass).then(({ project, token }) => {
    history.push(`/mission-control/${project}`)
    localStorage.setItem("token", token)
    loadConfig(project)
  }).catch(({ error }) => {
    notification.error({
      message: 'Error while login',
      description: error
    });
  })
}

export const saveConfig = () => {
  const projectId = get(store.getState(), "config.id")
  const config = get(store.getState(), "config")
  const result = unAdjustConfig(config)
  if (!result.ack) {
    if (Object.keys(result.errors.crud).length) {
      let errorDesc = ''
      Object.keys(result.errors.crud).forEach(db => {
        errorDesc = `${errorDesc} ${db} - (${result.errors.crud[db].join(", ")})`
      })
      notification.error({
        message: 'Error in CRUD config',
        description: errorDesc,
        duration: 0
      });
    }

    if (result.errors.functions.length) {
      notification.error({
        message: 'Error in Functions config',
        description: `Services - ${result.errors.functions.join(", ")}`,
        duration: 0
      });
    }
    if (result.errors.fileStore.length) {
      notification.error({
        message: 'Error in File Storage config',
        description: `Rules - ${result.errors.fileStore.join(", ")}`,
        duration: 0
      });
    }
    return
  }
  service.saveConfig(projectId, result.config).then(() => {
    notification.success({
      message: 'Success',
      description: "Config saved successfully"
    });
    store.dispatch(set("savedConfig", store.getState().config))
  }).catch(({ error }) => {
    notification.error({
      message: 'Error saving config',
      description: error
    });
  })
}