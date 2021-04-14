#!/bin/bash

if [ "$1" == "" ]; then
    echo "fotoDen test environment script"
    echo -ne "\n"
    echo "build         : builds fotoDen in the current folder (requires this to be in the same directory as the source)"
    echo "mk-testdir    : creates the fotoDen test directory"
    echo "mk-container  : makes a NGINX container (requires docker and sudo)"
    echo "run-container : runs the fotoDen-test container with the fotoDen_test folder mounted (expects a nginx container, requires sudo)"
    echo "rm-container  : removes the fotoDen-test container (requires sudo)"
    echo "mkall         : makes everything"
    exit 0
fi

if [ "$1" == "build" ] || [ "$1" == "mkall" ]; then
    make all
fi

if [ "$1" == "mk-testdir" ] || [ "$1" == "mkall" ]; then
    mkdir fotoDen_test

    fotoDen=$PWD'/build/fotoDen --config-dir fotoDen_test/'

    # $fotoDen init js build/fotoDen.min.js
    unzip build/default_theme.zip -d build/theme/
    $fotoDen init theme build/theme/
    $fotoDen init site --name "Test Site" -v --url "http://localhost" --theme Default fotoDen_test/test_root
    $fotoDen generate folder --name "Test Folder" -v fotoDen_test/test_root/test_folder
    $fotoDen generate album  --name "Test Album" -v test_images fotoDen_test/test_root/test_folder/test_album

    echo "---------------------------------------------------"
    echo "Your test environment is available at fotoDen_test."
    echo "---------------------------------------------------"
fi

if [ "$1" == "run-container" ]; then
    sudo docker run --name fotoDen-test -v $PWD/fotoDen_test/test_root:/usr/share/nginx/html:ro -p 80:80 -d nginx:alpine
elif [ "$1" == "rm-container" ]; then
    sudo docker rm -f fotoDen-test
elif [ "$1" == "mk-container" ] || [ "$1" == "mkall" ]; then
    read -n 1 -p "Would you like to make a NGINX container? (requires sudo) [y] " choice
    echo ""

    if [ "$choice" == "y" ] || [ "$choice" == "Y" ]; then
        sudo docker run --name fotoDen-test -v $PWD/fotoDen_test/test_root:/usr/share/nginx/html:ro -p 80:80 -d nginx:alpine
        echo "Access the test environment on http://localhost:80. Stop the environment by running docker stop fotoDen-test."
    fi
fi
