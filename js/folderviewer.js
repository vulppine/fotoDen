import currentTheme from './theme.js'
import Viewer, { contentLoad } from './viewer.js'
import { getJSON, getFolderURL, setText, setTitle } from './globals.js'

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

export const folderLoad = new CustomEvent('fd-folderLoad', { bubbles: true })

export default class FolderViewer extends Viewer {
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
          this.folderLinks.appendChild(currentTheme.loaded.createFolderLink(json))
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
