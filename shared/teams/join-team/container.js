// @flow
import * as TeamsGen from '../../actions/teams-gen'
import JoinTeamDialog from '.'
import {upperFirst} from 'lodash-es'
import {connect, compose, lifecycle, withStateHandlers, withHandlers} from '../../util/container'
import {type RouteProps} from '../../route-tree/render-route'

type OwnProps = RouteProps<{}, {}>

const mapStateToProps = state => ({
  errorText: upperFirst(state.teams.teamJoinError),
  success: state.teams.teamJoinSuccess,
  successTeamName: state.teams.teamJoinSuccessTeamName,
})

const mapDispatchToProps = (dispatch, {navigateUp}: OwnProps) => ({
  _onJoinTeam: (teamname: string) => {
    dispatch(TeamsGen.createJoinTeam({teamname}))
  },
  _onSetTeamJoinError: (error: string) => {
    dispatch(TeamsGen.createSetTeamJoinError({error}))
  },
  _onSetTeamJoinSuccess: (success: boolean, teamname: string) => {
    dispatch(TeamsGen.createSetTeamJoinSuccess({success, teamname}))
  },
  onBack: () => dispatch(navigateUp()),
})

export default compose(
  connect<OwnProps, _, _, _, _>(
    mapStateToProps,
    mapDispatchToProps,
    (s, d, o) => ({...o, ...s, ...d})
  ),
  withStateHandlers({name: ''}, {onNameChange: () => (name: string) => ({name: name.toLowerCase()})}),
  withHandlers({
    onSubmit: ({name, _onJoinTeam}) => () => _onJoinTeam(name),
  }),
  lifecycle({
    componentDidMount() {
      this.props._onSetTeamJoinError('')
      this.props._onSetTeamJoinSuccess(false, null)
    },
  })
)(JoinTeamDialog)
