#!/bin/bash

kubectl apply -f svc_user.yaml
kubectl apply -f svc_apigw.yaml
kubectl apply -f svc_upload.yaml

