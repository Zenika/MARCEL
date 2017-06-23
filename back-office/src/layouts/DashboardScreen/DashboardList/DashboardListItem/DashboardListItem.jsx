//@flow
import React from 'react'
import CardMedia from 'react-toolbox/lib/card/CardMedia'
import CardTitle from 'react-toolbox/lib/card/CardTitle'
import CardText from 'react-toolbox/lib/card/CardText'
import CardActions from 'react-toolbox/lib/card/CardActions'
import Button from 'react-toolbox/lib/button/Button'
import DashboardCard from '../DashboardCard'
import type { Dashboard } from '../../../../dashboard/type'

import './DashboardListItem.css'

export type PropsType = {
  dashboard: Dashboard,
  selectDashboard: () => void,
  deleteDashboard: () => void,
}

const DashboardListItem = (props: PropsType) => {
  const { dashboard, selectDashboard, deleteDashboard } = props
  return (
    <DashboardCard>
      <CardMedia
        aspectRatio="wide"
        image="https://placeimg.com/800/450/nature"
      />
      <CardTitle title={dashboard.name} />
      <CardText>{dashboard.description}</CardText>
      <CardActions className="buttons">
        <Button label="modifier" onClick={selectDashboard} />
        <Button label="ouvrir" />
        <Button label="supprimer" />
      </CardActions>
    </DashboardCard>
  )
}

export default DashboardListItem