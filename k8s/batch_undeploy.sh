#!/bin/bash

kubectl delete -f svc_user.yaml
kubectl delete -f svc_apigw.yaml
kubectl delete -f svc_upload.yaml
