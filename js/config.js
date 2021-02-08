import { getTheme } from './theme.js'
import { BaseURL, getJSON } from './globals.js'

const config = {
  websiteTitle: null,
  workingDirectory: null,
  storageURLBase: null,
  imageRootDir: null,
  thumbnailFrom: null,
  displayImageFrom: null,
  downloadSizes: null,
  imageSizes: new Map()
}

export function setConfig () {
  return getJSON(BaseURL + '/config.json')
    .then(async function (json) {
      await readConfig(json)
      if (json.Theme === true) {
        await getTheme()
      }
      return Promise.resolve()
    })
    .catch((err) => {
      console.error('An error occurred: ' + err)
      getJSON(new URL(document.URL).origin.toString())
        .then(json => {
          readConfig(json)
        })
        .catch(error => Promise.reject(error))
    })
}

export const configLoad = new Event('fd-configLoad', { bubbles: true })

export function readConfig (info) {
  config.websiteTitle = info.WebsiteTitle
  config.storageURLBase = info.PhotoURLBase
  config.thumbnailFrom = info.ThumbnailFrom
  config.imageRootDir = info.ImageRootDir
  config.downloadSizes = info.DownloadSizes

  const p = new URL(BaseURL).pathname
  if (p === '' || p === '/') {
    config.workingDirectory = ''
    console.log(config.workingDirectory)
  } else {
    const pa = p.split('/')
    console.log(pa)
    if (pa[pa.length - 1] === '') {
      pa.pop()
      config.workingDirectory = pa[pa.length - 1]
    } else {
      config.workingDirectory = pa[pa.length - 1]
    }
  }

  config.displayImageFrom = {
    size: info.DisplayImageFrom,
    prefix: info.DisplayImageFrom + '_'
  }

  if (info.DisplayImageFrom === 'src') {
    config.displayImageFrom.prefix = ''
  }

  info.ImageSizes.forEach((i) => {
    config.imageSizes.set(i.SizeName, {
      directory: [config.imageRootDir, i.Directory].join('/'),
      prefix: i.SizeName + '_',
      localBool: i.LocalBool
    })
  })

  document.dispatchEvent(configLoad)
}

export { config as default }
