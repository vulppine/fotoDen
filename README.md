fotoDen
=======

*A statically structured, front-end dynamic photo gallery*

[https://github.com/vulppine/fotoDen/wiki](Wiki)

Using
-----

fotoDen is designed so that it can be used independent of any separate tool, but
it works well if you have the fotoDen tool in order to create folders and albums
that are directly accessable by fotoDen's front end.

Using the pre-built fotoDen tool from release is as simple as initializing a
website (with an optional theme - one is already included!)

``` sh
fotoDen init -i site --name "[ your website name here ]" my_website/
```

From there, you can start generating folders and albums as you wish:

``` sh
fotoDen create folder --name "[ your folder name here ]"
my_website/my_folder/ fotoDen create album --name "[ your album name here ]"
my/images/are/here my_website/my_folder/my_album
```

Run the fotoDen command for more options. (More detailed information and
commands will be added soon, including use of the build system!)

### Warning

fotoDen is still in its *very early* stages, at v0 - everything and anything is
subject to change. While I'll attempt to keep its base structure stable,
commands and folder metadata are *subject to change in drastic ways*. Upgrade
paths will be offered for fotoDen version transitions that up the minor version
number (v0 . *y* . 0), so that version transitions with new changes are more
easier.

Installing
----------

You will require the latest version of
[https://github.com/libvips/libvips](libvips) in order to use fotoDen.

You can either install from an existing build (which includes fotoDen.js md5
checksums for both the minified version, and the JS file that came with its
release), or you can install the tool alone by running ~go install
github.com/vulppine/fotoDen~.

Building
--------

You will require the latest version of
[https://github.com/libvips/libvips](libvips), [https://terser.org](terser), and
[https://golang.org](Go) in order to build fotoDen from source.

1. Clone the git repository into a directory of your choice
2. Run ~make all~ in the resulting folder
3. fotoDen will be located in the **build/** folder in the same directory.

Contributing
------------

If you want to contribute, I encourage you to fork and help develop fotoDen!
Note that fotoDen uses [https://github.com/standard/standard](Standard JS) for
its JavaScript style.

A testing script is included in the root of the repository - run it in your
local environment to generate a test website. Generating the container will
require Docker.

Dependencies
------------

fotoDen relies on a few important dependencies (both front and back), so here
are links to the dependencies that fotoDen uses! License information for
included source code can be found in their respective locations.

| Dependency         | Author                 | License    |
|--------------------|------------------------|------------|
| [bimg]             | Tomas Aparicio         | MIT        |
| [Cobra]            | spf13                  | Apache 2.0 |
| [Bootstrap]        | Twitter/Bootstrap Team | MIT        |
| [justified-layout] | Flickr/SmugMug         | MIT        |
| [exif-js]          | Jacob Seidelin         | MIT        |
| [go-yaml]          | go-yaml team           | Apache 2.0 |
| [Goldmark]         | goldmark team          | MIT        |

[bimg]: https://github.com/h2non/bimg
[Cobra]: https://github.com/spf13/cobra
[Bootstrap]: https://github.com/twbs/bootstrap
[justified-layout]: https://github.com/flickr/justified-layout
[exif-js]: https://github.com/exif-js/exif-js
[go-yaml]: https://github.com/go-yaml/yaml
[Goldmark]: https://github.com/yuin/goldmark


Copyright
---------

fotoDen is copyright 2021 Flipp Syder under the MIT License (see LICENSE for
more information)

All test images licensed are under the
[https://creativecommons.org/licenses/by-nc-sa/4.0/](CC-BY-NC-SA)
