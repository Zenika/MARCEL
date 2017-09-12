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

class DashboardListItem extends React.Component {
  props: {
    dashboard: Dashboard,
    frontendURI: string,
    selectDashboard: () => void,
    deleteDashboard: () => void,
  }

  onClickWithoutPropagation = (onClick: () => void) => (event: MouseEvent) => {
    event.stopPropagation()
    onClick()
  }

  selectDashboard = this.onClickWithoutPropagation(this.props.selectDashboard)
  deleteDashboard = this.onClickWithoutPropagation(this.props.deleteDashboard)
  openDashboard = this.onClickWithoutPropagation(() => {
    const { frontendURI, dashboard } = this.props
    window.open(frontendURI + dashboard.id)
    window.focus()
  })

  render() {
    const { dashboard } = this.props
    const { selectDashboard, deleteDashboard, openDashboard } = this
    return (
      <DashboardCard onClick={selectDashboard}>
        <CardMedia aspectRatio="wide" image="https://placeimg.com/800/450/nature" />
        <CardTitle title={dashboard.name} />
        <CardText>{dashboard.description}</CardText>
        <CardActions className="buttons">
          <Button icon="mode_edit" label="modifier" onClick={selectDashboard} />
          <Button icon="exit_to_app" label="ouvrir" onClick={openDashboard} />
          <Button icon="delete" label="supprimer" onClick={deleteDashboard} />
        </CardActions>
      </DashboardCard>
    )
  }
}

export default DashboardListItem
