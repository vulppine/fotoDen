import currentTheme from './theme.js'
import config from './config.js'
import Viewer, { contentLoad } from './viewer.js'
import { getPageInfo, getAlbumURL, getJSON, makePhotoURL, setLink, setText, setTitle, setCurrentURLParam } from './globals.js'

function PhotoObject () {
  this.name = ''
  this.desc = ''
  this.index = parseInt(getPageInfo(new URL(document.URL)).index)
}

export default class PhotoViewer extends Viewer {
  constructor (container, info) {
    super(container, info)
    const photo = new PhotoObject()

    if (isNaN(parseInt(photo.index))) {
      photo.index = 0
    }

    getJSON('itemsInfo.json')
      .then(json => {
        if (json.Metadata === true) {
          getJSON(config.imageRootDir + '/meta/' + json.ItemsInFolder[photo.index] + '.json')
            .then(meta => {
              if (meta.ImageName === '') {
                photo.name = json.ItemsInFolder[photo.index]
                photo.desc = 'No description provided...'
              } else {
                photo.name = meta.ImageName
                photo.desc = json.ImageDesc
              }

              setText(this.name, photo.name)
              setText(this.desc, photo.desc)
              setTitle([photo.name, photo.album])
            })
        } else {
          photo.name = json.ItemsInFolder[photo.index]
          setText(this.name, photo.name)
          setTitle([photo.name, photo.album])
        }

        photo.album = this.info.FolderName

        if (parseInt(photo.index) === json.ItemsInFolder.length - 1) {
          currentTheme.loaded.setButton(this.navPrev, setCurrentURLParam('index', (parseInt(photo.index) - 1)))
          currentTheme.loaded.setButton(this.navNext)
        } else if (parseInt(photo.index) === 0) {
          currentTheme.loaded.setButton(this.navPrev)
          currentTheme.loaded.setButton(this.navNext, setCurrentURLParam('index', 1))
        } else {
          currentTheme.loaded.setButton(this.navPrev, setCurrentURLParam('index', (parseInt(photo.index) - 1)))
          currentTheme.loaded.setButton(this.navNext, setCurrentURLParam('index', (parseInt(photo.index) + 1)))
        }

        setText(this.folderName, photo.album)
        setLink(this.folderName, getAlbumURL().toString())
        this.setPhoto(json.ItemsInFolder[photo.index])

        if (this.infoButtons !== null) {
          this.setDownloads(json.ItemsInFolder[photo.index])
        }
      })
  }

  setPhoto (image) {
    setText(name, image)
    this.container.querySelector('.fd-photo').src = makePhotoURL(
      config.displayImageFrom.prefix + image,
      config.imageSizes.get(config.displayImageFrom.size).directory,
      config.imageSizes.get(config.displayImageFrom.size).localBool
    )
    this.container.querySelector('.fd-photo').addEventListener('load', e => {
      this.container.querySelector('.fd-photo').dispatchEvent(contentLoad)
    })
  }

  setDownloads (image) {
    config.downloadSizes.forEach((value) => {
      const newButton = currentTheme.loaded.createButton(
        value,
        makePhotoURL(
          config.imageSizes.get(value).prefix + image,
          config.imageSizes.get(value).directory,
          config.imageSizes.get(value).localBool)
      )
      this.infoButtons.appendChild(newButton)
    })
  }
}
