//@flow
import React from 'react'
import Button from 'react-toolbox/lib/button/Button'
import Grid from '../Grid'
import { values } from 'lodash'

import { ActivationButton, OpenButton, DeleteDashboardButton } from '../../common'
import type { Dashboard as DashboardT } from '../type'

import './Dashboard.css'

export type PropsType = {
  dashboard: DashboardT,
  uploadLayout: () => void,
}

const Dashboard = (props: PropsType) => {
  const { dashboard, uploadLayout } = props
  const { name, description, rows, cols, ratio, plugins } = dashboard
  return (
    <div className="Dashboard">
      <div className="head">
        <h2>
          {name} <br />
          <small>{description}</small>
        </h2>
        <div className="actions">
          <Button label="Sauvegarder" icon="save" primary onClick={uploadLayout} />
          <OpenButton dashboard={dashboard} />
          <DeleteDashboardButton dashboard={dashboard} />
          <ActivationButton dashboard={dashboard} />
        </div>
      </div>
      <Grid
        ratio={ratio}
        rows={rows}
        cols={cols}
        layout={values(plugins).map(({ x, y, cols, rows, ...instance }) => ({
          layout: { x, y, h: rows, w: cols },
          plugin: instance,
        }))}
      />
    </div>
  )
}

export default Dashboard
