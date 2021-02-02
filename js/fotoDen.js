/* eslint-env browser */

// fotoDen v0.1.0
//
// The front-end for a photo gallery.

/**
 * @version v0.1.0
 * @license MIT
 *
 * Copyright (c) 2021 Flipp Syder
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

const BaseURL = document.getElementById('fd-script').dataset.fdBaseurl
const version = '0.1.0-DEVEL'

// global variables

let isMobile

// configuration

let websiteTitle
let workingDirectory
let storageURLBase
let imageRootDir
let thumbnailFrom
let displayImageFrom
let downloadSizes
const imageSizes = new Map()

// theme

let theme = {
  config: {
    navRange: 5,
    imagesPerPage: 50
  },

  setButton: (button, text) => {
    if (button === null) { return }

    if (URL === undefined) {
      button.removeAttribute('href')
      button.setAttribute('class', 'button null')
    } else {
      button.href = URL
    }
  },

  setError: (errText) => {
    document.getElementsByClassName('errorBox')[0].setAttribute('style', 'display: block')
    document.getElementsByClassName('error')[0].innerHTML = errText
  },

  createThumbnail: (photoIndex, photoName) => {
    const thumbnailContainer = document.createElement('div')
    const thumbnail = new Image()
    const thumbnailAnchor = document.createElement('a')
    const thumbnailLink = new URL(document.URL)
    const thumbnailLinkParams = new URLSearchParams(thumbnailLink)

    thumbnailLinkParams.set('index', photoIndex)
    thumbnailLink.pathname = getAlbumURL().pathname.split('/').slice(0, getAlbumURL().pathname.split('/').length - 1).concat(['photo.html']).join('/')
    thumbnailLink.search = thumbnailLinkParams.toString()

    thumbnailContainer.appendChild(thumbnailAnchor)
    thumbnailContainer.setAttribute('class', 'fd-albumThumbnail')

    thumbnailAnchor.appendChild(thumbnail)
    thumbnailAnchor.setAttribute('href', thumbnailLink.toString())
    thumbnailAnchor.setAttribute('class', 'fd-albumThumbnailLink')

    thumbnail.setAttribute('class', 'albumThumbnailImage')
    thumbnail.setAttribute('src', makePhotoURL(imageSizes.get(thumbnailFrom).prefix + photoName, imageSizes.get(thumbnailFrom).directory, imageSizes.get(thumbnailFrom).localBool))

    return thumbnailContainer
  },

  createFolderLink: (info) => {
    const folderLinkContainer = document.createElement('div')
    const folderLink = getFolderURL(0).toString() + info.FolderShortName + '/'
    const folderInfoContainer = document.createElement('div')
    const folderItemCount = document.createElement('div')
    const folderAnchor = document.createElement('a')
    const folderThumbnail = new Image()

    folderAnchor.setAttribute('class', 'folderLink')
    folderLinkContainer.setAttribute('class', 'folderLinkContainer')
    folderInfoContainer.setAttribute('class', 'folderInfoContainer')
    folderItemCount.setAttribute('class', 'folderItemCount')
    folderThumbnail.setAttribute('class', 'folderThumbnail')

    folderAnchor.appendChild(folderLinkContainer)
    folderLinkContainer.appendChild(folderThumbnail)
    folderLinkContainer.appendChild(folderInfoContainer)
    folderInfoContainer.appendChild(folderItemCount)

    if (info.FolderThumbnail === true) {
      folderThumbnail.src = info.FolderShortName + '/' + info.FolderThumbnail
    } else {
      folderThumbnail.src = getAlbumURL() + info.FolderShortName + '/thumb.png'
    };

    if (info.ItemAmount != null) {
      const newDiv = document.createElement('div')
      newDiv.innerHTML = 'Photos: ' + info.ItemAmount // remember, this is still photo oriented...
      folderItemCount.appendChild(newDiv)
    };

    if (info.SubfolderShortNames.length > 0) {
      const newDiv = document.createElement('div')
      newDiv.innerHTML = 'Folders: ' + info.SubfolderShortNames.length
      folderItemCount.appendChild(newDiv)
    }

    const name = document.createElement('span')
    name.innerHTML = info.FolderName
    folderInfoContainer.insertBefore(name, folderItemCount)
    folderAnchor.href = folderLink

    return folderAnchor
  },

  createNavPageLink (page, container) {
    const newAnchor = document.createElement('a')
    const newURL = getAlbumURL()

    newAnchor.innerHTML = (page + 1)

    if (page === this.currentPage) {
      newAnchor.setAttribute('class', 'navLink active')
    } else {
      newURL.search = '?page=' + page
      newAnchor.href = newURL.toString()
      newAnchor.setAttribute('class', 'navLink')
    }

    container.appendChild(newAnchor)
  }
}

// generic functions

function setCurrentURLParam (param, value) {
  const newURL = new URL(document.URL)
  const newURLParams = new URLSearchParams(newURL.search)

  newURLParams.set(param, value)
  newURL.search = newURLParams.toString()

  return newURL
}

function getJSON (url) {
  if (!isMobile) {
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

function setTitle (items) {
  items.push(websiteTitle)
  document.title = items.join(' - ')
}

function setText (element, text) {
  if (element === null) { return }

  element.innerText = text
}

function setLink (element, link) {
  if (element === null) { return }

  element.href = link
}

function getFolderURL (level) {
  const folderURL = new URL(document.URL)
  const folderPath = folderURL.pathname.split('/').slice(0, folderURL.pathname.split('/').length - 1) // knock off any index.htmls or nulls right off the bat;

  let rootDirectoryLoc

  if (workingDirectory === '') {
    if (folderURL.pathname === '/' && level > 0) {
      return null
    }
    rootDirectoryLoc = 0
  } else {
    rootDirectoryLoc = folderPath.indexOf(workingDirectory)
  }

  folderURL.search = ''

  if (rootDirectoryLoc !== 0 && folderPath.length - level < folderPath.length - rootDirectoryLoc) {
    console.warn('Attempted to go deeper than baseURL, ignoring.')
    return null
  } else if (level <= folderPath.length) {
    folderURL.pathname = folderPath.slice(0, folderPath.length - level).concat(['']).join('/') // folders should really have a default page file name
    folderURL.href = folderURL.origin + folderURL.pathname + folderURL.search // had an issue with this, so i'm forcing it
    return folderURL
  } else {
    console.warn('Attempted to go deeper than baseURL, ignoring.')
    return null
  }
}

function getPageInfo (url) {
  const search = new URLSearchParams(url.search)

  return {
    index: search.get('index'),
    page: search.get('page')
  }
}

// photoViewer

/* photoObject
 *
 * The main idea is that the photo.html page should be
 * a static page that's linked from the generated
 * album pages, if somebody wants a bigger view of a photo
 * but they also want to link directly back to the site
 *
 * Most of the design of the display is already within
 * the HTML file, this lirary's task is to just give it
 * the actual functionality it requires.
 *
 */

function makePhotoURL (photoName, dir, localBool) {
  if (storageURLBase === 'local' || storageURLBase === '' || localBool === true) {
    return getAlbumURL() + dir + '/' + photoName
  } else {
    const newURL = new URL(document.URL)
    const newURLPathArray = newURL.pathname.split('/')
    const rootDirectoryLoc = newURLPathArray.indexOf(workingDirectory)

    return storageURLBase + newURLPathArray.slice(rootDirectoryLoc, newURLPathArray.length - 1).join('/') + '/' + dir + '/' + photoName
  }
}

function PhotoObject () {
  this.name = ''
  this.desc = ''
  this.index = parseInt(getPageInfo(new URL(document.URL)).index)
}

// Viewer
//
// The base class for all other viewer types
// Construct it with folderInfo.json already called.
// This will handle:
// - album information
// - subfolder information
// - navigation setting
// - etc.

const viewerLoad = new Event('fd-viewerLoad', { bubbles: true })
const contentLoad = new Event('fd-contentLoad', { bubbles: true })

class Viewer {
  constructor (container, info) {
    this.container = container

    // CONFIG/CURRENT LOADED INFORMATION //
    this.info = info
    this.navContentRange = theme.config.navRange

    // SPANS/LINKS/TEXT //
    this.name = container.querySelector('.fd-name')
    this.desc = container.querySelector('.fd-desc')
    this.folderName = container.querySelector('.fd-folderName')
    this.superFolder = container.querySelector('.fd-superFolder')

    // CONTAINERS //
    this.folderSubtitle = container.querySelector('.fd-folderSubtitle')
    this.infoButtons = container.querySelector('.fd-infoButtons')
    this.navContents = container.querySelector('.fd-navContents')

    // BUTTONS //
    this.navNext = container.querySelector('.fd-navNext')
    this.navPrev = container.querySelector('.fd-navPrev')

    setText(this.folderName, info.FolderName)
    this.setSuperFolder()

    container.dispatchEvent(viewerLoad)
  }

  setSuperFolder () {
    const f = getFolderURL(1)
    if (f === null) {
      if (this.folderSubtitle !== null) {
        this.folderSubtitle.setAttribute('style', 'display: none')
      }
    } else {
      getJSON(f.toString() + 'folderInfo.json')
        .then((info) => {
          setText(this.superFolder, info.FolderName)
          setLink(this.superFolder, f.toString())
        })
    }
  }

  getNavContentMinMax (total, current) {
    if (total < this.navContentRange) {
      return [0, total]
    }

    let eachSide
    if (this.navContentRange % 2 !== 0) {
      eachSide = (this.navContentRange - 1) / 2
    }

    if (current - eachSide < 0) {
      return [0, this.navContentRange]
    } else if (current + eachSide > total) {
      return [total - this.navContentRange, total]
    } else {
      return [current - eachSide, current + eachSide + 1]
    }
  }
}

const PhotoViewers = []

class PhotoViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    const photo = new PhotoObject()

    if (isNaN(parseInt(photo.index))) {
      photo.index = 0
    }

    getJSON('itemsInfo.json')
      .then(json => {
        if (json.Metadata === true) {
          getJSON(imageRootDir + '/meta/' + json.ItemsInFolder[photo.index] + '.json')
            .then(meta => {
              if (meta.ImageName === '') {
                photo.name = json.ItemsInFolder[photo.index]
              } else {
                photo.name = meta.ImageName
              }

              photo.desc = json.ImageDesc
              setText(this.name, photo.name)
              setText(this.desc, photo.desc)
            })
        } else {
          photo.name = json.ItemsInFolder[photo.index]
          setText(this.name, photo.name)
        }

        photo.album = this.info.FolderName

        if (parseInt(photo.index) === json.ItemsInFolder.length - 1) {
          theme.setButton(this.navPrev, setCurrentURLParam('index', (parseInt(photo.index) - 1)))
          theme.setButton(this.navNext)
        } else if (parseInt(photo.index) === 0) {
          theme.setButton(this.navPrev)
          theme.setButton(this.navNext, setCurrentURLParam('index', 1))
        } else {
          theme.setButton(this.navPrev, setCurrentURLParam('index', (parseInt(photo.index) - 1)))
          theme.setButton(this.navNext, setCurrentURLParam('index', (parseInt(photo.index) + 1)))
        }

        setText(this.folderName, photo.album)
        setLink(this.folderName, getAlbumURL().toString())
        setTitle([photo.name, photo.album])
        this.setPhoto(json.ItemsInFolder[photo.index])

        if (this.infoButtons !== null) {
          this.setDownloads(json.ItemsInFolder[photo.index])
        }
      })
  }

  setPhoto (image) {
    setText(name, image)
    this.container.querySelector('.fd-photo').src = makePhotoURL(
      displayImageFrom.prefix + image,
      imageSizes.get(displayImageFrom.size).directory,
      imageSizes.get(displayImageFrom.size).localBool
    )
    this.container.querySelector('.fd-photo').addEventListener('load', e => {
      this.container.querySelector('.fd-photo').dispatchEvent(contentLoad)
    })
  }

  setDownloads (image) {
    downloadSizes.forEach((value) => {
      const newButton = theme.createButton(
        value,
        makePhotoURL(
          imageSizes.get(value).prefix + image,
          imageSizes.get(value).directory,
          imageSizes.get(value).localBool)
      )
      this.infoButtons.appendChild(newButton)
    })
  }
}

// albums

function getAlbumURL () {
  return getFolderURL(0)
}

// AlbumViewer
//
// AlbumViewers are containers for album-type pages.
// They take in a specific JSON object, and use its info in order to display the album
// within the current folder (specified by the URL of the folder).
//
// Functions within an AlbumViewer include generating navigation pages,
// and calling the thumbnail generator to create thumbnails according
// to the current items in the folder, from ImageThumbDir from the current website's config.

const AlbumViewers = []

const imageLoad = new Event('fd-imgLoad', { bubbles: true })

// constructor for AlbumViewer

class AlbumViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    this.imagesPerPage = theme.config.imagesPerPage
    this.currentPage = parseInt(getPageInfo(new URL(document.URL)).page)

    if (isNaN(this.currentPage)) { this.currentPage = 0 }

    this.photos = null
    this.maxPhotos = null
    this.pageAmount = null

    this.thumbnailContainer = this.container.querySelector('.fd-albumThumbnails')
    this.currentThumbnails = null

    if (this.container.querySelector('.fd-folder.fd-viewer') !== null) {
      this.folderViewer = this.container.querySelector('.fd-folder.fd-viewer')
    };

    setTitle([info.FolderName])
    if (info.FolderDesc !== '') {
      setText(this.desc, info.FolderDesc)
    }

    getJSON(getAlbumURL() + 'itemsInfo.json')
      .then((json) => {
        this.photos = json.ItemsInFolder
        this.maxPhotos = this.photos.length
        this.pageAmount = Math.ceil(this.maxPhotos / this.imagesPerPage)

        this.update()
      })
      .catch(error => {
        console.error(error)
        theme.setError('Could not load album properly. Code: ' + error)
      }) // you can't throw out of a constructor, sadly, so we'll just make it visible to the user

    if (isMobile) {
      this.thumbnailContainer.addEventListener('scroll', () => {
        const currentScroll = this.thumbnailContainer.scrollTop
        const maxHeight = this.thumbnailContainer.scrollHeight

        if (currentScroll > (maxHeight - (maxHeight * 0.25))) {
          if (this.currentPage !== Math.ceil(this.maxPhotos / this.imagesPerPage)) {
            this.currentPage++
            this.populate()
          }
        }
      })
    }
  }

  setNavPageLinks () {
    if (this.currentPage === Math.ceil(this.maxPhotos / this.imagesPerPage) - 1) {
      theme.setButton(this.navPrev)
      theme.setButton(this.navNext)
    } else if (this.currentPage === 0) {
      theme.setButton(this.navPrev)
      theme.setButton(this.navNext, setCurrentURLParam('page', (this.currentPage + 1)))
    } else if (this.currentPage === this.pageAmount) {
      theme.setButton(this.navPrev, setCurrentURLParam('page', (this.currentPage - 1)))
      theme.setButton(this.navNext)
    } else {
      theme.setButton(this.navPrev, setCurrentURLParam('page', (this.currentPage - 1)))
      theme.setButton(this.navNext, setCurrentURLParam('page', (this.currentPage + 1)))
    }

    const range = this.getNavContentMinMax(Math.ceil(this.maxPhotos / this.imagesPerPage), this.currentPage)

    for (let i = range[0]; i < range[1]; i++) {
      console.log(this.navContents)
      theme.createNavPageLink(i, this.navContents)
    }
  }

  populate () {
    let index = this.imagesPerPage * this.currentPage

    while (index < this.maxPhotos) {
      if (index === (this.imagesPerPage * this.currentPage) + this.imagesPerPage) {
        break
      } else {
        const newThumbnail = theme.createThumbnail(index, this.photos[index])
        newThumbnail.getElementsByTagName('img')[0].addEventListener('load', () => newThumbnail.dispatchEvent(imageLoad))
        this.thumbnailContainer.appendChild(newThumbnail)
      }
      index++
    }

    let totalLoaded = 0
    this.thumbnailContainer.addEventListener('fd-imgLoad', () => {
      totalLoaded++

      if (totalLoaded === this.imagesPerPage || totalLoaded === this.info.ItemAmount) {
        this.thumbnailContainer.dispatchEvent(contentLoad)
      }
    })
  }

  update () {
    this.currentThumbnails = this.container.getElementsByClassName('albumThumbnail')
    for (let i = 0; i < this.currentThumbnails.length; i++) {
      this.currentThumbnails.item(0).remove()
    };

    if (this.folderViewer !== undefined) {
      if (this.currentPage !== 0) {
        this.folderViewer.remove()
      };
    };

    this.setNavPageLinks()
    this.populate()
  }
}

/* folderViewer
 *
 * Generic sub-folder viewer.
 *
 */

// Constructor for a FolderViewer
//
// Allows access to a container that contains a FolderViewer.
// Includes access for setting any current titles, compared to an AlbumViewer which auto-sets.
// If the current FolderType is 'folder', it will auto-set the style to a full version on init instead of the shorter, smaller version seen in albums.

const FolderViewers = []

const folderLoad = new CustomEvent('fd-folderLoad', { bubbles: true })

class FolderViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    this.folderLinks = container.querySelector('.fd-folderLinks')
    this.style = ''

    this.type = info.FolderType

    if (this.info.FolderType !== 'album') {
      setTitle([info.FolderName])
      if (info.FolderDesc !== '') {
        setText(this.desc, info.FolderDesc)
      }
    }

    if (this.info.SubfolderShortNames.length > 0) {
      this.populate()
    } else {
      this.container.remove()
    }
  }

  populate () {
    this.info.SubfolderShortNames.forEach(element => {
      getJSON(getFolderURL(0).toString() + element + '/folderInfo.json')
        .then(json => {
          this.folderLinks.appendChild(theme.createFolderLink(json))
          this.container.dispatchEvent(folderLoad)
        })
    })

    let totalLoaded = 0
    this.container.addEventListener('fd-folderLoad', () => {
      totalLoaded++
      console.log(totalLoaded)

      if (totalLoaded === this.info.SubfolderShortNames.length) {
        this.folderLinks.dispatchEvent(contentLoad)
      }
    })
  }
}

function setConfig () {
  return getJSON(BaseURL + '/config.json')
    .then(async function (json) {
      readConfig(json)
      if (json.Theme === true) {
        theme = await import(BaseURL + '/theme/js/theme.js')
        theme.init()
      }
      return Promise.resolve()
    })
    .catch(() => getJSON(new URL(document.URL).origin.toString())
      .then(json => {
        readConfig(json)
      })
      .catch(error => Promise.reject(error)))
}

function readConfig (info) {
  websiteTitle = info.WebsiteTitle
  storageURLBase = info.PhotoURLBase
  thumbnailFrom = info.ThumbnailFrom
  imageRootDir = info.ImageRootDir
  downloadSizes = info.DownloadSizes

  const p = new URL(BaseURL).pathname
  if (p === '' || p === '/') {
    workingDirectory = ''
    console.log(workingDirectory)
  } else {
    const pa = p.split('/')
    console.log(pa)
    if (pa[pa.length] === '') {
      pa.pop()
      workingDirectory = pa[pa.length]
    } else {
      workingDirectory = pa[pa.length]
    }
  }

  displayImageFrom = {
    size: info.DisplayImageFrom,
    prefix: info.DisplayImageFrom + '_'
  }

  if (info.DisplayImageFrom === 'src') {
    displayImageFrom.prefix = ''
  }

  info.ImageSizes.forEach((i) => {
    imageSizes.set(i.SizeName, {
      directory: [imageRootDir, i.Directory].join('/'),
      prefix: i.SizeName + '_',
      localBool: i.LocalBool
    })
  })
}

/* Page initilization
 *
 * Checks if we're in a photo viewer,
 * otherwise attempts to initialize both the folders and the album items.
 */

function pageInit () {
  if (BaseURL === null) {
    throw new Error('Error: fd-baseURL was not defined in the HTML. Aborting.')
  }

  setText(document.getElementById('fd-version'), version)
  const mobileCheck = window.matchMedia('(pointer: coarse)')
  if (mobileCheck.matches) { isMobile = true }
  setConfig()
    .then(() => {
      getJSON(getAlbumURL() + 'folderInfo.json')
        .then((info) => {
          try {
            const viewers = document.querySelectorAll('.fd-viewer')
            viewers.forEach((viewer) => {
              if (viewer.classList.contains('fd-photo')) {
                const photoViewer = new PhotoViewer(viewer, info)
                PhotoViewers.push(photoViewer)
              } else if (viewer.classList.contains('fd-folder')) {
                const folderViewer = new FolderViewer(viewer, info)
                FolderViewers.push(folderViewer)
              } else if (viewer.classList.contains('fd-album')) {
                const albumViewer = new AlbumViewer(viewer, info)
                AlbumViewers.push(albumViewer)
              } else {
                console.warn('Invalid viewer type detected.')
              }
            })
          } catch (err) {
            console.error(err)
          }
        })
        .catch(error => {
          setTitle([error])
          console.error(error)
          theme.setError('Error getting folder information: ' + error)
        })
    })
    .catch((status) => {
      setTitle([status])
      console.error(status)
      theme.setError('error getting config: ' + status)
    })
}

window.onload = pageInit()
