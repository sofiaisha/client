// @flow
import * as React from 'react'
import * as RPCChatTypes from '../../constants/types/rpc-chat-gen'
import * as Sb from '../../stories/storybook'
import Chat from '.'
import {Box} from '../../common-adapters/index'

const props = {
  unfurlMode: RPCChatTypes.unfurlUnfurlMode.whitelisted,
  unfurlWhitelist: ['amazon.com', 'wsj.com', 'nytimes.com', 'keybase.io', 'google.com', 'twitter.com'],
  onUnfurlSave: (mode: RPCChatTypes.UnfurlMode, whitelist: Array<string>) => {
    Sb.action('onUnfurlSave')(mode, whitelist)
  },
}

const load = () => {
  Sb.storiesOf('Settings/Chat', module)
    .addDecorator(story => <Box style={{padding: 5}}>{story()}</Box>)
    .add('Default', () => <Chat {...props} />)
}

export default load
