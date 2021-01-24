/* eslint-env browser */

// fotoDen v0.0.1
//
// The front-end for a photo gallery.
//
// Copyright (c) 2021 Flipp Syder
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
//  in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// global variables

let isMobile

// configuration

let websiteTitle
let workingDirectory
let storageURLBase
let imageRootDir
let thumbnailFrom
let displayImageFrom
const imageSizes = new Map()

// theme
// note: this is polyfill until themes are implemented fully

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
    thumbnailContainer.setAttribute('class', 'albumThumbnail')

    thumbnailAnchor.appendChild(thumbnail)
    thumbnailAnchor.setAttribute('href', thumbnailLink.toString())
    thumbnailAnchor.setAttribute('class', 'albumThumbnailLink')

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

    if (info.ItemsInFolder != null) {
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

function getFolderURL (level) {
  const folderURL = new URL(document.URL)
  const folderPath = folderURL.pathname.split('/').slice(0, folderURL.pathname.split('/').length - 1) // knock off any index.htmls or nulls right off the bat;

  let rootDirectoryLoc

  if (workingDirectory === '') {
    rootDirectoryLoc = 0
  } else {
    rootDirectoryLoc = folderPath.indexOf(workingDirectory)
  }

  folderURL.search = ''

  if (rootDirectoryLoc !== 0 && folderPath.length - level < folderPath.length - rootDirectoryLoc) {
    throw new Error('Error: attempted to go deeper than workingDirectory, aborting!')
  } else if (level <= folderPath.length) {
    folderURL.pathname = folderPath.slice(0, folderPath.length - level).concat(['']).join('/') // folders should really have a default page file name
    folderURL.href = folderURL.origin + folderURL.pathname + folderURL.search // had an issue with this, so i'm forcing it
    return folderURL
  } else {
    throw new Error('Attempted to go deeper than possible - ignoring.')
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
  this.album = '' // we get this from folderInfo.json
  this.index = getPageInfo(new URL(document.URL)).index
  this.canvas = null
  this.blobURL = null
  this.cloudURL = null
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

class Viewer {
  constructor (container, info) {
    this.container = container

    // CONFIG/CURRENT LOADED INFORMATION //
    this.info = info
    this.navContentRange = theme.config.navRange

    // SPANS/LINKS/TEXT //
    this.name = container.querySelector('.name')
    this.folderName = container.querySelector('.folderName')
    this.superFolder = container.querySelector('.superFolder')

    // CONTAINERS //
    this.folderSubtitle = container.querySelector('.folderSubtitle')
    this.infoButtons = container.querySelector('.infoButtons')
    this.navContents = container.querySelector('.navContents')

    // BUTTONS //
    this.navNext = container.querySelector('.navNext')
    this.navPrev = container.querySelector('.navPrev')
  }

  setSuperFolder () {
    getJSON(getFolderURL(1).toString() + 'folderInfo.json')
      .then((info) => {
        this.superFolder.innerHTML = info.FolderName
        this.superFolder.href = getFolderURL(1).toString()
      })
      .catch(() => {
        if (this.folderSubtitle === null) { return }

        this.folderSubtitle.setAttribute('style', 'display: none')
      })
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
    photo.name = this.info.ItemsInFolder[photo.index]
    photo.album = this.info.FolderName
    this.setAlbum(photo.album)

    if (parseInt(photo.index) === this.info.ItemsInFolder.length - 1) {
      this.setPrev(setCurrentURLParam('index', (parseInt(photo.index) - 1)))
      this.setNext(null)
    } else if (parseInt(photo.index) === 0) {
      this.setPrev(null)
      this.setNext(setCurrentURLParam('index', 1))
    } else {
      this.setPrev(setCurrentURLParam('index', (parseInt(photo.index) - 1)))
      this.setNext(setCurrentURLParam('index', (parseInt(photo.index) + 1)))
    }

    setTitle([photo.name, photo.album])
    this.setPhoto(this.info.ItemsInFolder[photo.index])
  }

  setAlbum (name) {
    this.folderName.innerHTML = name
    this.folderName.href = getAlbumURL().toString()
  }

  setPhoto (image) {
    setText(name, image)
    document.getElementsByClassName('mainPhoto')[0].src = makePhotoURL(displayImageFrom.prefix + image, imageSizes.get(displayImageFrom.size).directory, imageSizes.get(displayImageFrom.size).localBool)
  }

  setDownloads (image) {
    imageSizes.forEach((value, key) => {
      const newButton = document.createElement('a')
      newButton.setAttribute('class', 'downloadButton button')

      if (key === 'src') {
        newButton.innerHTML = 'src'
        newButton.href = makePhotoURL(image, value.directory, value.localBool)
      } else {
        newButton.innerHTML = key
        newButton.href = makePhotoURL(key + '_' + image, value.directory, value.localBool)
      }
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

    this.thumbnailContainer = this.container.querySelector('.thumbnailContainer')
    this.currentThumbnails = null

    if (this.container.getElementsByClassName('folderLinks').length === 1) {
      this.folderViewer = this.container.getElementsByClassName('folderLinks')[0]
    };

    setTitle([info.FolderName])
    this.setSuperFolder()

    getJSON(getAlbumURL() + 'itemsInfo.json')
      .then((json) => {
        this.photos = json.ItemsInFolder
        this.maxPhotos = this.photos.length
        this.pageAmount = Math.ceil(this.maxPhotos / this.imagesPerPage)

        this.update()
      })
      .catch(error => { theme.setError('Could not load album properly. Code: ' + error) }) // you can't throw out of a constructor, sadly, so we'll just make it visible to the user

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

  createNavPageLink (page) {
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

    return newAnchor
  }

  setNavPageLinks () {
    if (this.currentPage === 1) {
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
      this.navContents.appendChild(this.createNavPageLink(i))
    }
  }

  populate () {
    let index = this.imagesPerPage * this.currentPage

    while (index < this.maxPhotos) {
      if (index === (this.imagesPerPage * this.currentPage) + this.imagesPerPage) {
        break
      } else {
        this.thumbnailContainer.appendChild(theme.createThumbnail(index, this.photos[index]))
      }
      index++
    }
  }

  update () {
    this.currentThumbnails = this.container.getElementsByClassName('albumThumbnail')
    for (let i = 0; i < this.currentThumbnails.length; i++) {
      this.currentThumbnails.item(0).remove()
    };

    if (this.folderViewer !== undefined) {
      if (this.currentPage === 1 && this.info.SubfolderShortNames[0] !== undefined) {
        // todo: put something here?
      } else {
        this.folderViewer.setAttribute('style', 'display: none')
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

class FolderViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    this.folderLinks = container.querySelector('.folderLinks')
    this.style = ''

    this.type = info.FolderType

    if (this.info.FolderType !== 'album') {
      setText(this.name, info.FolderName)
      setTitle([info.FolderName])
      this.setSuperFolder()
    }

    this.populate()
  }

  populate () {
    this.info.SubfolderShortNames.forEach(element => {
      getJSON(getFolderURL(0).toString() + element + '/folderInfo.json')
        .then(json => this.folderLinks.appendChild(theme.createFolderLink(json)))
    })
  }
}

function setConfig () {
  return getJSON('/config.json')
    .then(json => {
      if (json.Theme !== '') {
        import('/js/theme.js')
          .then(extheme => {
            theme = extheme
          })
      }
      readConfig(json)
      return Promise.resolve()
    })
    .catch(() => getJSON(new URL(document.URL).origin.toString())
      .then(json => {
        readConfig(json)
        return Promise.resolve()
      })
      .catch(error => Promise.reject(error)))
}

function readConfig (info) {
  websiteTitle = info.WebsiteTitle
  workingDirectory = info.WorkingDirectory
  storageURLBase = info.PhotoURLBase
  thumbnailFrom = info.ThumbnailFrom
  imageRootDir = info.ImageRootDir

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
  const mobileCheck = window.matchMedia('(pointer: coarse)')
  if (mobileCheck.matches) { isMobile = true }
  setConfig()
    .then(() => {
      getJSON(getAlbumURL() + 'folderInfo.json')
        .then((info) => {
          try {
            if (document.getElementById('PhotoViewer')) {
              const photoViewer = new PhotoViewer(document.getElementById('PhotoViewer'), info)
              PhotoViewers.push(photoViewer)
            } else {
              if (document.getElementById('FolderViewer')) {
                const folderViewer = new FolderViewer(document.getElementById('FolderViewer'), info)
                FolderViewers.push(folderViewer)
              }
              if (document.getElementById('AlbumViewer')) {
                const albumViewer = new AlbumViewer(document.getElementById('AlbumViewer'), info)
                AlbumViewers.push(albumViewer)
              }
            };
          } catch (err) {
            console.log(err)
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
