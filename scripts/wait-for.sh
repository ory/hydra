#!/bin/bash

# Change Port value to Hydra server
while ! nc -z localhost 4444; do   
  sleep 10
done
