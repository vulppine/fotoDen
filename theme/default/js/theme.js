/* global bootstrap, BaseURL, getPageInfo, getAlbumURL, makePhotoURL, imageSizes, thumbnailFrom, websiteTitle */
/* eslint-env browser */
// fotoDen DEFAULT THEME FUNCTIONS //

// internal functions //

function toggleView (container) {
  if (container.classList.contains('d-none')) {
    container.classList.remove('d-none')
  } else {
    container.classList.add('d-none')
  }
}

function removeLoad (container) {
  container.getElementsByClassName('fd-loading')[0].remove()
}

const imageRatios = []
const loadedImages = []

function getLoadedImageRatios (container) {
  const thumbnails = container.getElementsByClassName('fd-albumThumbnail')
  for (let i = 0; i < thumbnails.length; i++) {
    loadedImages.push(thumbnails[i].getElementsByTagName('img')[0])
  }

  loadedImages.forEach(img => {
    imageRatios.push({
      width: img.width,
      height: img.height
    })
  })
}

async function justifyThumbnails (container) {
  const layout = await import(BaseURL + '/theme/js/layout.js')

  const newImageSizes = layout.justifyImages(imageRatios, { containerWidth: container.offsetWidth })

  newImageSizes.boxes.forEach((newSize, i) => {
    loadedImages[i].width = newSize.width
    loadedImages[i].height = newSize.height
    loadedImages[i].setAttribute('style', 'padding: 5px;')
  })

  window.onresize = () => {
    justifyThumbnails(container)
  }
}

// exported functions //

export const config = {
  navRange: 5,
  imagesPerPage: 50
}

export function setError (errText) {
  const errorBox = new bootstrap.Modal(document.getElementsByClassName('fd-error')[0])

  document.querySelector('.errorText').innerText = errText
  errorBox.toggle()
}

export function setButton (button, URL) {
  if (button === null) { return }

  // pagination specific //

  if (button.classList.contains('page-item') === true) {
    console.log('setting nav')
    if (URL === undefined) {
      button.classList.add('disabled')
      button.querySelector('.page-link').classList.remove('text-white')
      button.querySelector('.page-link').removeAttribute('href')
      return
    } else {
      button.querySelector('.page-link').setAttribute('href', URL)
      return
    }
  }

  if (URL === undefined) {
    button.removeAttribute('href')
    if (button.classList.contains('disabled') !== true) {
      button.classList.add('disabled')
    }
  } else {
    button.setAttribute('href', URL)
    button.classList.remove('disabled')
  }
}

export function createButton (name, URL) {
  const anchor = document.createElement('a')
  anchor.href = URL
  const button = document.createElement('button')
  anchor.appendChild(button)

  button.setAttribute('class', 'btn btn-secondary mx-1')
  button.innerHTML = name

  return anchor
}

export function createNavPageLink (page, container) {
  const newPage = document.createElement('li')
  const newAnchor = document.createElement('a')
  newPage.appendChild(newAnchor)
  newPage.classList.add('page-item')
  newAnchor.classList.add('page-link', 'border-white')
  const newURL = getAlbumURL()

  let currentPage = parseInt(getPageInfo(new URL(document.URL)).page)
  if (isNaN(currentPage)) { currentPage = 0 }

  newAnchor.innerHTML = (page + 1)

  if (page === currentPage) {
    newPage.classList.add('active')
    newAnchor.classList.add('bg-white', 'text-dark')
  } else {
    newURL.search = '?page=' + page
    newAnchor.href = newURL.toString()
    newAnchor.classList.add('bg-dark', 'text-white')
  }

  container.insertBefore(newPage, container.querySelector('.fd-navNext'))
}

export function createFolderLink (info) {
  const folderContainer = document.createElement('div')
  const folderAnchor = document.createElement('a')
  const folderCard = document.createElement('div')
  const folderThumb = document.createElement('img')
  const folderInfoContainer = document.createElement('div')
  const folderInfoGrid = document.createElement('div')
  const folderTitleContainer = document.createElement('div')
  const folderTitle = document.createElement('h5')
  const folderItemCountContainer = document.createElement('div')

  folderContainer.appendChild(folderAnchor)
  folderContainer.setAttribute('class', 'col-lg h-75 fd-folderLink')

  folderAnchor.appendChild(folderCard)
  folderAnchor.setAttribute('class', 'text-white text-decoration-none')
  folderAnchor.href = info.FolderShortName

  folderCard.appendChild(folderThumb)
  folderCard.appendChild(folderInfoContainer)
  folderCard.setAttribute('class', 'card bg-dark overflow-hidden')

  folderThumb.setAttribute('class', 'card-img')
  folderThumb.setAttribute('style', 'height: 250px; object-fit: cover')

  folderInfoContainer.appendChild(folderInfoGrid)
  folderInfoContainer.setAttribute('class', 'card-img-overlay d-flex align-items-end p-0')

  folderInfoGrid.appendChild(folderTitleContainer)
  folderInfoGrid.appendChild(folderItemCountContainer)
  folderInfoGrid.setAttribute('class', 'col bg-dark p-2')

  folderTitleContainer.appendChild(folderTitle)
  folderTitleContainer.setAttribute('class', 'row')

  if (info.FolderThumbnail === true) {
    folderThumb.src = info.FolderShortName + '/thumb.jpg'
  } else {
    folderThumb.classList.add('bg-secondary')
  }

  folderItemCountContainer.setAttribute('class', 'row')

  if (info.ItemAmount > 0) {
    const folderItemCountPhotos = document.createElement('div')
    folderItemCountContainer.appendChild(folderItemCountPhotos)
    folderItemCountPhotos.setAttribute('class', 'col')

    const folderItemCountPhotosText = document.createElement('small')
    folderItemCountPhotos.appendChild(folderItemCountPhotosText)
    folderItemCountPhotosText.setAttribute('class', 'mb-0 card-text')

    folderItemCountPhotosText.innerHTML = '<i class="bi bi-camera-fill"></i> : ' + info.ItemAmount
  }

  if (info.SubfolderShortNames.length > 0) {
    const folderItemCountFolders = document.createElement('div')
    folderItemCountContainer.appendChild(folderItemCountFolders)
    folderItemCountFolders.setAttribute('class', 'col')

    const folderItemCountFolderText = document.createElement('small')
    folderItemCountFolders.appendChild(folderItemCountFolderText)
    folderItemCountFolderText.setAttribute('class', 'mb-0 card-text')

    folderItemCountFolderText.innerHTML = '<i class="bi bi-folder-fill"></i> : ' + info.SubfolderShortNames.length
  }

  folderTitle.innerText = info.FolderName

  return folderContainer
}

export function createThumbnail (index, name) {
  const thumbnail = new Image()
  const thumbnailAnchor = document.createElement('a')
  const thumbnailLink = new URL(document.URL)
  const thumbnailLinkParams = new URLSearchParams(thumbnailLink)

  thumbnailLinkParams.set('index', index)
  thumbnailLink.pathname = getAlbumURL().pathname.split('/').slice(0, getAlbumURL().pathname.split('/').length - 1).concat(['photo.html']).join('/')
  thumbnailLink.search = thumbnailLinkParams.toString()

  thumbnail.setAttribute('class', 'fd-albumThumbnailImage')
  thumbnail.setAttribute('src', makePhotoURL(imageSizes.get(thumbnailFrom).prefix + name, imageSizes.get(thumbnailFrom).directory, imageSizes.get(thumbnailFrom).localBool))

  thumbnailAnchor.appendChild(thumbnail)
  thumbnailAnchor.href = thumbnailLink.toString()

  thumbnailAnchor.classList.add('fd-albumThumbnail', 'text-center')

  return thumbnailAnchor
}

export function init () {
  document.addEventListener('fd-viewerLoad', e => {
    const t = e.target
    if (t.classList.contains('fd-root')) {
      removeLoad(document.querySelector('#fd-rootLoad'))
      toggleView(t)
      if (e.target.classList.contains('fd-album') || e.target.classList.contains('fd-folder')) {
        console.log('Root album/folder container detected')
        fetch('./thumb.jpg')
          .then(response => {
            if (!response.ok) {
              console.error(response)
              throw new Error()
            }
            return response.blob()
          })
          .then((blob) => {
            console.log('Thumbnail detected')
            t.querySelector('.fd-banner').getElementsByTagName('img')[0].src = URL.createObjectURL(blob)
            t.querySelector('.fd-bannerload').remove()
          })
          .catch((err) => {
            console.log(e.target)
            console.log(err)
            t.querySelector('.fd-banner').remove()
          })
      }
    }

    console.log(e.target)
  })

  document.addEventListener('fd-contentLoad', e => {
    console.log(e.target)
    if (e.target.classList.contains('fd-albumThumbnails')) {
      getLoadedImageRatios(e.target)
      justifyThumbnails(e.target)
    }

    removeLoad(e.target.parentNode)
    toggleView(e.target)
  })

  if (document.querySelector('.navbar-brand') !== null) {
    if (websiteTitle === null || websiteTitle === undefined || websiteTitle === '') {
      document.querySelector('.navbar-brand').innerText = 'fotoDen'
    } else {
      document.querySelector('.navbar-brand').innerText = websiteTitle
    }
  }
}
