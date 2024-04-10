#!/bin/bash
# Deploy etcd inspections

set -o errexit

env=
num=
cluster_ip=
kubectl create ns kstone
kubectl -n kstone create sa kstone
kubectl apply -f clusterrolebindings.yaml
kubectl create secret generic ${env}-etcd -n kstone --from-file=/etc/kubernetes/pki/etcd/ca.crt --from-file=/etc/kubernetes/pki/etcd/server.crt --from-file=/etc/kubernetes/pki/etcd/server.key
sed -i "s/env/${env}/g" etcd-cluster.yaml
sed -i "s/num/${num}/g" etcd-cluster.yaml
sed -i "s/cluster_ip/${cluster_ip}/g" etcd-cluster.yaml
kubectl apply -f etcd-cluster.yaml
kubectl apply -f etcdclusters-deployment.yaml
kubectl apply -f etcdinspections-deployment.yaml
kubectl apply -f etcdinspections-svc.yaml
kubectl apply -f servicemonitor.yaml