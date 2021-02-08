import config from './config.js'

// global variables

export const BaseURL = document.getElementById('fd-script').dataset.fdBaseurl

// global functions

export function isMobile () {
  return window.matchMedia('(pointer: coarse)').matches
}

export function debug (func) {
  if (document.getElementById('fd-script').dataset.fdDebug === 'true') {
    return func()
  }
}

export function makePhotoURL (photoName, dir, localBool) {
  if (config.storageURLBase === 'local' || config.storageURLBase === '' || localBool === true) {
    return getAlbumURL() + dir + '/' + photoName
  } else {
    const newURL = new URL(document.URL)
    const newURLPathArray = newURL.pathname.split('/')
    const rootDirectoryLoc = newURLPathArray.indexOf(config.workingDirectory)

    return config.storageURLBase + newURLPathArray.slice(rootDirectoryLoc, newURLPathArray.length - 1).join('/') + '/' + dir + '/' + photoName
  }
}

export function setCurrentURLParam (param, value) {
  const newURL = new URL(document.URL)
  const newURLParams = new URLSearchParams(newURL.search)

  newURLParams.set(param, value)
  newURL.search = newURLParams.toString()

  return newURL
}

export function getJSON (url) {
  if (!window.isMobile) {
    return fetch(url)
      .then(response => {
        if (!response.ok) {
          throw new Error(response.status)
        }
        return Promise.resolve(response.json())
      })
  } else {
    const request = new XMLHttpRequest()

    return new Promise((resolve, reject) => {
      request.onreadystatechange = () => {
        if (request.readyState === 4) {
          if (request.status === 200) {
            resolve(JSON.parse(request.response))
          } else {
            reject(request.status)
          }
        }
      }
      request.open('GET', url)
      request.send()
    })
  }
}

// setTitle
//
// Takes an array of strings, separates them via a separator string and adds
// the website title at the end.

export function setTitle (items) {
  items.push(config.websiteTitle)
  document.title = items.join(' - ')
}

export function setText (element, text) {
  if (element === null) { return }

  element.innerText = text
}

export function setLink (element, link) {
  if (element === null) { return }

  element.href = link
}

export function getAlbumURL () {
  return getFolderURL(0)
}

export function getFolderURL (level) {
  const folderURL = new URL(document.URL)
  const folderPath = folderURL.pathname.split('/').slice(0, folderURL.pathname.split('/').length - 1) // knock off any index.htmls or nulls right off the bat;

  let rootDirectoryLoc

  if (config.workingDirectory === '') {
    if (folderURL.pathname === '/' && level > 0) {
      return null
    }
    rootDirectoryLoc = 0
  } else {
    rootDirectoryLoc = folderPath.indexOf(config.workingDirectory)
  }

  folderURL.search = ''

  if (rootDirectoryLoc !== 0 && folderPath.length - 1 - level < rootDirectoryLoc) {
    debug(console.warn('Attempted to go deeper than baseURL, ignoring.'))
    return null
  } else if (level <= folderPath.length) {
    folderURL.pathname = folderPath.slice(0, folderPath.length - level).concat(['']).join('/') // folders should really have a default page file name
    folderURL.href = folderURL.origin + folderURL.pathname + folderURL.search // had an issue with this, so i'm forcing it
    return folderURL
  } else {
    debug(console.warn('Attempted to go deeper than baseURL, ignoring.'))
    return null
  }
}

export function getPageInfo (url) {
  const search = new URLSearchParams(url.search)

  return {
    index: search.get('index'),
    page: search.get('page')
  }
}
