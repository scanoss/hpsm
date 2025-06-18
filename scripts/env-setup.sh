#!/bin/bash

##########################################
#
# This script will copy all the required files into the correct locations on the server
# Binary goes into: /usr/local/bin
# Library goes into: /usr/local/lib
#
################################################################

if [ "$1" = "-h" ] || [ "$1" = "-help" ] ; then
  echo "$0 [-help] [-f]"
  echo "   Setup and copy the relevant files into place on a server to run the SCANOSS HPSM Feature"
  echo "   -f    force installation (accept all prompts)"
  exit 1
fi

export B_PATH=/usr/local/bin/
export L_PATH=/usr/local/lib/

# Check for -f flag
ACCEPT_ALL=false
if [ "$1" = "-f" ] ; then
  ACCEPT_ALL=true
fi

# Makes sure the scanoss user exists
export RUNTIME_USER=scanoss
if ! getent passwd $RUNTIME_USER > /dev/null ; then
  echo "Runtime user does not exist: $RUNTIME_USER"
  echo "Please create using: useradd --system $RUNTIME_USER"
  exit 1
fi
# Also, make sure we're running as root
if [ "$EUID" -ne 0 ] ; then
  echo "Please run as root"
  exit 1
fi
if [ "$ACCEPT_ALL" = true ] ; then
  echo "Auto-accepting installation (non-interactive mode)..."
  echo "Starting installation..."
else
  read -p "Install SCANOSS HPSM (y/n) [n]? " -n 1 -r
  echo
  if [[ $REPLY =~ ^[Yy]$ ]] ; then
    echo "Starting installation..."
  else
    echo "Stopping."
    exit 1
  fi
fi

# Setup the service on the system (defaulting to service name without environment)

B_FILENAME="hpsm"
L_FILENAME="libhpsm.so"

echo "Copying HPSM binary..."
if [ -f "$B_FILENAME" ] ; then
  if ! cp "$B_FILENAME" "$B_PATH" ; then
    echo "binary copy failed"
    exit 1
  fi
fi
if [ -f "$L_FILENAME" ] ; then
  if ! cp "$L_FILENAME" "$L_PATH" ; then
    echo "library copy failed"
    exit 1
  fi
  if ! ldconfig ; then
    echo "updating ldconfig failed"
    exit 1
  fi
fi

echo
echo "Review binary in: $B_PATH"
echo "Review library in: $L_PATH"
echo