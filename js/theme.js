import config from './config.js'
import { BaseURL, getAlbumURL, getFolderURL, makePhotoURL } from './globals.js'

const currentTheme = {
  loaded: { // EXTREME JANK
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
      document.getElementsByClassName('fd-errorBox')[0].setAttribute('style', 'display: block')
      document.getElementsByClassName('fd-error')[0].innerHTML = errText
    },

    createButton: (name, url) => {
      const button = document.createElement('a')
      button.classList.Add('button')
      button.href = url
      button.innerText = name

      return button
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
      thumbnail.setAttribute('src', makePhotoURL(config.imageSizes.get(config.thumbnailFrom).prefix + photoName, config.imageSizes.get(config.thumbnailFrom).directory, config.imageSizes.get(config.thumbnailFrom).localBool))

      return thumbnailContainer
    },

    createFolderLink: (info) => {
      const folderLinkContainer = document.createElement('div')
      const folderLink = getFolderURL(0).toString() + info.FolderShortName + '/'
      const folderInfoContainer = document.createElement('div')
      const folderItemCount = document.createElement('div')
      const folderAnchor = document.createElement('a')
      const folderThumbnail = new Image()

      folderAnchor.setAttribute('class', 'fd-folderLink')
      folderLinkContainer.setAttribute('class', 'fd-folderLinkContainer')
      folderInfoContainer.setAttribute('class', 'fd-folderInfoContainer')
      folderItemCount.setAttribute('class', 'fd-folderItemCount')
      folderThumbnail.setAttribute('class', 'fd-folderThumbnail')

      folderAnchor.appendChild(folderLinkContainer)
      folderLinkContainer.appendChild(folderThumbnail)
      folderLinkContainer.appendChild(folderInfoContainer)
      folderInfoContainer.appendChild(folderItemCount)

      if (info.FolderThumbnail === true) {
        folderThumbnail.src = info.FolderShortName + '/' + 'thumb.jpg'
      } else {
        folderThumbnail.src = BaseURL + '/thumb.png'
      }

      if (info.ItemAmount != null) {
        const newDiv = document.createElement('div')
        newDiv.innerHTML = 'Photos: ' + info.ItemAmount // remember, this is still photo oriented...
        folderItemCount.appendChild(newDiv)
      }

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
        newAnchor.setAttribute('class', 'fd-navLink active')
      } else {
        newURL.search = '?page=' + page
        newAnchor.href = newURL.toString()
        newAnchor.setAttribute('class', 'fd-navLink')
      }

      container.appendChild(newAnchor)
    }
  }
}

export async function getTheme (url) {
  import(BaseURL + '/theme/js/theme.js')
    .then((t) => {
      currentTheme.loaded = t
      currentTheme.loaded.init()
      console.log(currentTheme.loaded)
      return Promise.resolve()
    })
    .catch((err) => {
      console.error('An error occurred, using default settings: ' + err)
      return Promise.resolve()
    }) 
}

export default currentTheme
