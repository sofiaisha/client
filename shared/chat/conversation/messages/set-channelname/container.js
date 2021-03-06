// @flow
import {connect, isMobile} from '../../../../util/container'
import * as ProfileGen from '../../../../actions/profile-gen'
import * as Types from '../../../../constants/types/chat2'
import * as TrackerGen from '../../../../actions/tracker-gen'
import SetChannelname from '.'

type OwnProps = {|
  message: Types.MessageSetChannelname,
|}

const mapStateToProps = (state, {message}) => ({
  author: message.author,
  channelname: message.newChannelname,
  setUsernameBlack: message.author === state.config.username,
  timestamp: message.timestamp,
})

const mapDispatchToProps = (dispatch, {message}) => ({
  onUsernameClicked: () =>
    isMobile
      ? dispatch(ProfileGen.createShowUserProfile({username: message.author}))
      : dispatch(
          TrackerGen.createGetProfile({forceDisplay: true, ignoreCache: true, username: message.author})
        ),
})

export default connect<OwnProps, _, _, _, _>(
  mapStateToProps,
  mapDispatchToProps,
  (s, d, o) => ({...o, ...s, ...d})
)(SetChannelname)
