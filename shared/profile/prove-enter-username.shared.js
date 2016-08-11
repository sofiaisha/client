// @flow
import type {PlatformsExpandedType} from '../constants/types/more'

function standardText (name: string) {
  return {
    headerText: `Prove your ${name} identity`,
    floatingLabelText: `Your ${name} username`,
    hintText: `Your ${name} username`,
  }
}

export const platformText : {[key: PlatformsExpandedType]: {headerText: string, floatingLabelText?: string, hintText?: string}} = {
  'twitter': standardText('Twitter'),
  'reddit': standardText('Reddit'),
  'github': standardText('GitHub'),
  'coinbase': standardText('Coinbase'),
  'hackernews': standardText('Hacker News'),
  'btc': {
    headerText: 'Set a Bitcoin address',
    floatingLabelText: 'Your Bitcoin address',
  },
  'dns': {
    headerText: 'Prove your domain',
    hintText: 'yourdomain.com',
  },
  'http': {
    headerText: 'Prove your website',
    hintText: 'whatever.yoursite.com',
  },
  'https': {
    headerText: 'Prove your website',
    hintText: 'whatever.yoursite.com',
  },
}
