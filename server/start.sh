#!/bin/sh

echo "starting Xvfb on :99"
Xvfb :99 &
./backend
