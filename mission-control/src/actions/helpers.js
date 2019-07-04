import { cloneDeep } from "lodash"
import { notification } from "antd"

export const adjustConfig = (config) => {
  // Adjust database rules
  if (config.modules && config.modules.crud) {
    Object.keys(config.modules.crud).forEach(db => {
      if (config.modules.crud[db].collections) {
        Object.keys(config.modules.crud[db].collections).forEach(col => {
          config.modules.crud[db].collections[col] = JSON.stringify(config.modules.crud[db].collections[col], null, 2)
        })
      }
    })
  }

  // Adjust function rules
  if (config.modules && config.modules.functions && config.modules.functions.rules) {
    Object.keys(config.modules.functions.rules).forEach(service => {
      config.modules.functions.rules[service] = JSON.stringify(config.modules.functions.rules[service], null, 2)
    })
  }

  // Adjust file storage rules
  if (config.modules && config.modules.fileStore && config.modules.fileStore.rules) {
    Object.keys(config.modules.fileStore.rules).forEach(rule => {
      config.modules.fileStore.rules[rule] = JSON.stringify(config.modules.fileStore.rules[rule], null, 2)
    })
  }

  return config
}

export const unAdjustConfig = (c) => {
  let config = cloneDeep(c)
  let result = { ack: true, errors: { crud: {}, functions: [], fileStore: [] } }
  // Unadjust database rules
  if (config.modules && config.modules.crud) {
    Object.keys(config.modules.crud).forEach(db => {
      if (config.modules.crud[db].collections) {
        Object.keys(config.modules.crud[db].collections).forEach(col => {
          try {
            config.modules.crud[db].collections[col] = JSON.parse(config.modules.crud[db].collections[col])
          } catch (error) {
            result.ack = false
            result.errors.crud[db] = result.errors.crud[db] ? [...result.errors.crud[db], col] : [col]
          }
        })
      }
    })
  }

  // Unadjust function rules
  if (config.modules && config.modules.functions && config.modules.functions.rules) {
    Object.keys(config.modules.functions.rules).forEach(service => {
      try {
        config.modules.functions.rules[service] = JSON.parse(config.modules.functions.rules[service])
      } catch (error) {
        result.ack = false
        result.errors.functions.push(service)
      }
    })
  }

  // Unadjust file storage rules
  if (config.modules && config.modules.fileStore && config.modules.fileStore.rules) {
    Object.keys(config.modules.fileStore.rules).forEach(rule => {
      try {
        config.modules.fileStore.rules[rule] = JSON.parse(config.modules.fileStore.rules[rule])
      } catch (error) {
        result.ack = false
        result.errors.functions.push(rule)
      }
    })
  }

  result.config = config
  return result
}

export const notify = (type, title, msg, duration) => {
  notification[type]({ message: title, description: msg, duration: duration });
}