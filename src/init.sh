#!/bin/sh

# default
USERNAME="delobUser"
PASSWORD="delobPassword"

for arg in "$@"; do
    case $arg in
        USERNAME=*)
            USERNAME="${arg#USERNAME=}"
            ;;
        PASSWORD=*)
            PASSWORD="${arg#PASSWORD=}"
            ;;
        *)
            echo "Warning: Unrecognized argument: $arg"
            ;;
    esac
done

./elo --add-user $USERNAME $PASSWORD

./delob