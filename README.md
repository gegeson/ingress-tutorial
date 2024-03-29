# ingress-tutorial

## 手順

### 準備

以下、$VM の箇所は VM の IP に置き換える。
次のように調べられる。このケースでは 192.168.64.3 がホストから VM に接続するための IP となる。

    $multipass ls
    Name                    State             IPv4             Image
    microk8s-vm             Running           192.168.64.3     Ubuntu 18.04 LTS
                                            10.1.254.64

container image の build

    docker build -t $VM:32000/sample:latest .

container image の local registry への push
※事前に Docker Engine の設定に insecure-registries として追加する必要がある。

    {
    "insecure-registries" : ["$VM:32000"]
    }

    docker push $VM:32000/sample:latest

manifest の反映

    kubectl apply -f namespace.yaml
    kubectl apply -f deployment.yaml
    kubectl apply -f service.yaml
    kubectl apply -f ingress.yaml

デフォルト namespace の変更

    kubectl config set-context $(kubectl config current-context) --namespace=ingress-tutorial

リソースの確認

    kubectl get deploy,svc,ingress

動作確認

    open http://$VM

ログ出力

    kubectl exec sample-deployment-84d7b4f486-qsnr2 -- printenv

pathによってだしわけ

    $ kubectl describe ingress
    Name:             sample-ingress
    Namespace:        ingress-tutorial
    Address:          127.0.0.1
    Default backend:  sample-service:api (10.1.254.94:8000,10.1.254.95:8000,10.1.254.96:8000 + 1 more...)
    Rules:
    Host        Path  Backends
    ----        ----  --------
    *
                /      sample-service:8000 (10.1.254.94:8000,10.1.254.95:8000,10.1.254.96:8000 + 1 more...)
                /foo   sample-service2:8000 (10.1.254.94:8000,10.1.254.95:8000,10.1.254.96:8000 + 1 more...)
    Annotations:  <none>
    Events:
    Type    Reason  Age                  From                      Message
    ----    ------  ----                 ----                      -------
    Normal  Sync    12s (x4 over 3h35m)  nginx-ingress-controller  Scheduled for sync

参考
[Kubernetes API Reference Docs](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#envvar-v1-core)
