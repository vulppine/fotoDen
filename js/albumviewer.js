import currentTheme from './theme.js'
import Viewer, { contentLoad } from './viewer.js'
import { debug, getAlbumURL, getJSON, getPageInfo, isMobile, setCurrentURLParam, setText, setTitle } from './globals.js'

export const imageLoad = new Event('fd-imgLoad', { bubbles: true })

// constructor for AlbumViewer

export default class AlbumViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    this.imagesPerPage = currentTheme.loaded.config.imagesPerPage
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
        currentTheme.loaded.setError('Could not load album properly. Code: ' + error)
      }) // you can't throw out of a constructor, sadly, so we'll just make it visible to the user

    if (isMobile()) {
      if (this.navContents !== null) {
        this.navContents.remove()
      }

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
      currentTheme.loaded.setButton(this.navPrev)
      currentTheme.loaded.setButton(this.navNext)
    } else if (this.currentPage === 0) {
      currentTheme.loaded.setButton(this.navPrev)
      currentTheme.loaded.setButton(this.navNext, setCurrentURLParam('page', (this.currentPage + 1)))
    } else if (this.currentPage === this.pageAmount) {
      currentTheme.loaded.setButton(this.navPrev, setCurrentURLParam('page', (this.currentPage - 1)))
      currentTheme.loaded.setButton(this.navNext)
    } else {
      currentTheme.loaded.setButton(this.navPrev, setCurrentURLParam('page', (this.currentPage - 1)))
      currentTheme.loaded.setButton(this.navNext, setCurrentURLParam('page', (this.currentPage + 1)))
    }

    const range = this.getNavContentMinMax(Math.ceil(this.maxPhotos / this.imagesPerPage), this.currentPage)

    for (let i = range[0]; i < range[1]; i++) {
      debug(console.log(this.navContents))
      currentTheme.loaded.createNavPageLink(i, this.navContents)
    }
  }

  populate () {
    let index = this.imagesPerPage * this.currentPage

    while (index < this.maxPhotos) {
      if (index === (this.imagesPerPage * this.currentPage) + this.imagesPerPage) {
        break
      } else {
        const newThumbnail = currentTheme.loaded.createThumbnail(index, this.photos[index])
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
