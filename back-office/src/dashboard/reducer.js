//@flow
import type { Reducer } from 'redux'
import { actions } from './actions'
import { mapValues } from 'lodash'
import { set, update, unset } from 'immutadot'
import uuid from 'uuid/v4'
import type {
  DashboardAction,
  DashboardState,
  LayoutMap,
  PluginInstanceMap,
} from './type'

import mockedData from '../mocked-data/dashboards'
const intialState = mockedData

const updatePlugins = (layout: LayoutMap) => (plugins: PluginInstanceMap) => {
  return mapValues(plugins, plugin => {
    if (!layout[plugin.instanceId])
      throw new Error('Plugin instance not found in layout')

    const { x, y, w: columns, h: rows } = layout[plugin.instanceId]
    return { ...plugin, x, y, columns, rows }
  })
}

const dashboard: Reducer<DashboardState, DashboardAction> = (
  state = intialState,
  action,
) => {
  switch (action.type) {
    case actions.SELECT_PLUGIN: {
      return { ...state, selectedPlugin: action.payload.instanceId }
    }
    case actions.SELECT_DASHBOARD: {
      return { ...state, selectedDashboard: action.payload.dashboardName }
    }
    case actions.UNSELECT_DASHBOARD: {
      return { ...state, selectedDashboard: null, selectedPlugin: null }
    }
    case actions.ADD_PLUGIN: {
      const instanceId = uuid()
      const { selectedDashboard } = state
      return selectedDashboard
        ? set(state, `dashboards.${selectedDashboard}.plugins.${instanceId}`, {
            ...action.payload.plugin,
            x: 0,
            y: 0,
            columns: 1,
            rows: 1,
            instanceId,
          })
        : state
    }
    case actions.DELETE_PLUGIN: {
      const { selectedPlugin, selectedDashboard } = state
      return selectedPlugin
        ? selectedDashboard
          ? unset(
              state,
              `dashboards.${selectedDashboard}.plugins.${selectedPlugin}`,
            )
          : state
        : state
    }
    case actions.CHANGE_PROP: {
      const { instanceId, prop, value } = action.payload
      const { selectedDashboard } = state
      return selectedDashboard
        ? set(
            state,
            `dashboards.${selectedDashboard}.plugins.${instanceId}.props.${prop.name}.value`,
            value,
          )
        : state
    }
    case actions.SAVE_LAYOUT: {
      const { layout } = action.payload
      const { selectedDashboard } = state
      return selectedDashboard
        ? update(
            state,
            `dashboards.${selectedDashboard}.plugins`,
            updatePlugins(layout),
          )
        : state
    }
    case actions.UPDATE_CONFIG: {
      const { selectedDashboard } = state
      const { property, value } = action.payload
      const parsedValue = !isNaN(value) ? parseFloat(value) : value
      return selectedDashboard
        ? set(state, `dashboards.${selectedDashboard}.${property}`, parsedValue)
        : state
    }
    default:
      return state
  }
}

export default dashboard
