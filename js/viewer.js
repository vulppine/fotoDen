import currentTheme from './theme.js'
import { setText, getFolderURL, getJSON, setLink } from './globals.js'

export const viewerLoad = new Event('fd-viewerLoad', { bubbles: true })
export const contentLoad = new Event('fd-contentLoad', { bubbles: true })

export default class Viewer {
  constructor (container, info) {
    this.container = container

    // CONFIG/CURRENT LOADED INFORMATION //
    this.info = info
    this.navContentRange = currentTheme.loaded.config.navRange

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
