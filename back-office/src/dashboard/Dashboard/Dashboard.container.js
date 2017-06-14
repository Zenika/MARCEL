//@flow
import { connect } from 'react-redux'
import { dashboardSelector } from '../selectors'
import { deletePlugin } from '../actions'
import Dashboard from './Dashboard'

const mapStateToProps = state => ({
  dashboard: dashboardSelector(state),
})

const mapDispatchToProps = {
  deletePlugin,
}

export default connect(mapStateToProps, mapDispatchToProps)(Dashboard)
