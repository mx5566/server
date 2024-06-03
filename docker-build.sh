#!/bin/bash

sudo docker buildx build -t serveimage:1.0.0 . --platform linux/amd64 -f Dockerfile
