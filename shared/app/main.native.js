// @flow
import * as React from 'react'
import Push from './push.native'
import RenderRoute from '../route-tree/render-route'
import {connect, type TypedState} from '../util/container'
import {navigateUp, setRouteState} from '../actions/route-tree'
import {GatewayDest} from 'react-gateway'
import {NativeBackHandler} from '../common-adapters/mobile.native'
import {View} from 'react-native'
import {globalStyles} from '../styles'
import {isAndroid} from '../constants/platform'
import {getPath} from '../route-tree'

type Props = {
  folderBadge: number,
  routeDef: any,
  routeState: any,
  showPushPrompt: any,
  setRouteState: (path: any, partialState: any) => void,
  navigateUp: () => void,
}

class Main extends React.Component<Props> {
  componentDidMount() {
    if (isAndroid) {
      NativeBackHandler.addEventListener('hardwareBackPress', () => {
        if (getPath(this.props.routeState).size === 1) {
          return false
        }
        this.props.navigateUp()
        return true
      })
    }
  }

  render() {
    if (this.props.showPushPrompt) {
      return <Push />
    }

    return (
      <React.Fragment>
        <RenderRoute
          routeDef={this.props.routeDef}
          routeState={this.props.routeState}
          setRouteState={this.props.setRouteState}
        />
        <GatewayDest
          name="popup-root"
          component={View}
          pointerEvents="box-none"
          style={globalStyles.fillAbsolute}
        />
      </React.Fragment>
    )
  }
}

const mapStateToProps = (state: TypedState) => ({
  folderBadge: state.favorite.folderState.privateBadge + state.favorite.folderState.publicBadge,
  routeDef: state.routeTree.routeDef,
  routeState: state.routeTree.routeState,
  showPushPrompt: state.push.showPushPrompt,
})

const mapDispatchToProps = (dispatch: Dispatch) => ({
  navigateUp: () => dispatch(navigateUp()),
  setRouteState: (path, partialState) => dispatch(setRouteState(path, partialState)),
})

export default connect(mapStateToProps, mapDispatchToProps)(Main)