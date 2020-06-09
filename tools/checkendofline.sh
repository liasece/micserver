#!/bin/bash

find . -name "*.go" | xargs grep -r -n $'\r'
