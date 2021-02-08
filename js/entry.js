import currentTheme from './theme.js'
import { BaseURL, getAlbumURL, getJSON, setText, setTitle } from './globals.js'
import PhotoViewer from './photoviewer.js'
import FolderViewer from './folderviewer.js'
import AlbumViewer from './albumviewer.js'
import { setConfig } from './config.js'

export * from './globals.js'
export { default as config } from './config.js'

const version = '0.1.1'

const PhotoViewers = []
const AlbumViewers = []
const FolderViewers = []

function pageInit () {
  if (BaseURL === null) {
    throw new Error('Error: fd-baseURL was not defined in the HTML. Aborting.')
  }

  setText(document.getElementById('fd-version'), version)
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
          currentTheme.loaded.setError('Error getting folder information: ' + error)
        })
    })
    .catch((status) => {
      setTitle([status])
      console.error(status)
      currentTheme.loaded.setError('error getting config: ' + status)
    })
}

pageInit()
