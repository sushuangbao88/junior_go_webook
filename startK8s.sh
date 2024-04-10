echo "启动K8s \n"

kubectl apply -f webook-deployment.yaml
kubectl apply -f webook-service.yaml

kubectl apply -f webook-record-mysql-deployment.yaml
kubectl apply -f webook-record-mysql-pv.yaml
kubectl apply -f webook-record-mysql-pvc.yaml
kubectl apply -f webook-record-mysql-service.yaml

kubectl apply -f webook-record-redis-deployment.yaml
kubectl apply -f webook-record-redis-service.yaml

echo "启动完成! \n"

echo "查看deployment： "
kubectl get deployment|grep "webook"
echo "\n"

echo "查看service： "
kubectl get service|grep "webook"
echo "\n"

echo "查看pods： "
kubectl get pods|grep "webook"
echo "\n"

echo "查看PersistentVolume： "
kubectl get PersistentVolume|grep "webook"
echo "\n"

echo "查看PersistentVolumeClaim： "
kubectl get PersistentVolumeClaim|grep "webook"
echo "\n"