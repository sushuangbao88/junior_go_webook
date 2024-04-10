echo "关闭K8s \n"

kubectl delete service webook-record
kubectl delete service webook-record-mysql
kubectl delete service webook-record-redis

kubectl delete deployment webook-record
kubectl delete deployment webook-record-mysql
kubectl delete deployment webook-record-redis

kubectl delete PersistentVolume webook-mysql-pvc
kubectl delete PersistentVolume webook-mysql-pv


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