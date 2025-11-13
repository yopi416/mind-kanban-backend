# mind-kanban-backend
マインドマップ×カンバンボードアプリケーションのバックエンド開発

## AWS 環境

### ビルド＆デプロイ

#### ビルド＆プッシュ

（ローカル）

```powershell
aws configure

# 変数

## powershell用(変更してないので別PJのもの)
$REGION  = 'ap-northeast-1'
$ACCOUNT = '726101441058'
$REPO    = 'minkan'
$IMAGE   = "${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${REPO}"
$TAG     = 'local'

## Ubuntu用
REGION="ap-northeast-1"
ACCOUNT="726101441058"
REPO="minkan"
IMAGE="${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${REPO}"
TAG="latest"


# ECR ログイン
#aws ecr get-login-password --region ${REGION} `
#| docker login --username AWS --password-stdin "${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com"

aws ecr get-login-password --region ${REGION}| docker login --username AWS --password-stdin "${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com"

# ビルド & プッシュ
docker build -t "${IMAGE}:${TAG}" .
docker push "${IMAGE}:${TAG}"

# latest タグも付与してプッシュ
docker tag "${IMAGE}:${TAG}" "${IMAGE}:latest"
docker push "${IMAGE}:latest"
```

#### イメージプル sZA

- EC2 インスタンスに SSM から接続して以下のコマンドを実行する。



```sh
#######
# SSMで入ったとき専用 
sudo -i -u ubuntu
cd /home/ubuntu/project/minkan
# #####

aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 726101441058.dkr.ecr.ap-northeast-1.amazonaws.com

docker pull 726101441058.dkr.ecr.ap-northeast-1.amazonaws.com/minkan:latest
```

#### サーバー再起動

- EC2 インスタンスに SSM から接続して以下のコマンドを実行する。

```sh
CONTAINER_NAME=abe-tools

# 停止
docker rm -f "$CONTAINER_NAME"

# 起動
docker run -d --name $CONTAINER_NAME --network abe-net --restart unless-stopped --shm-size=1g --gpus all 339713090098.dkr.ecr.ap-northeast-1.amazonaws.com/abe-demo-blur:latest uvicorn main:app --host 0.0.0.0 --port 8081 --log-level info

# 確認
docker logs "$CONTAINER_NAME"
# [INFO] Application startup complete.
# が出ていれば成功。
```

### 鍵生成

#### マスター鍵

- 生成

```sh
cd ~
mkdir -p ./keys
docker run --rm --name abe-master-keygen --gpus all -v "$(pwd)/keys:/out" 339713090098.dkr.ecr.ap-northeast-1.amazonaws.com/abe-demo-blur:latest python3 pyztoolkit/kpfabeo_setup.py -p /out/

# → ~/keys に mpk.kpfabeo, msk.kpfabeo ができる
```

- Base64 に変換（DB 保存用）

```sh
base64 "$HOME/keys/mpk.kpfabeo" > "$HOME/keys/mpk.kpfabeo.b64"
base64 "$HOME/keys/msk.kpfabeo" > "$HOME/keys/msk.kpfabeo.b64"
```

#### ポリシー付き鍵

```sh
cd ~
mkdir -p ./out

docker run --rm --name abe-keygen --gpus all  -v "$(pwd)/keys:/keys:ro"   -v "$(pwd)/out:/out" 339713090098.dkr.ecr.ap-northeast-1.amazonaws.com/abe-demo-blur:latest bash -lc '
  set -euo pipefail
  POLICY="plate"
  OUTPUT_FILE="/out/plate.key"
  python3 pyztoolkit/kpfabeo_keygen.py -msk /keys/msk.kpfabeo -mpk /keys/mpk.kpfabeo -i "$POLICY" -o $OUTPUT_FILE
  echo "===== BEGIN $OUTPUT_FILE (Base64) ====="
  base64 $OUTPUT_FILE > $OUTPUT_FILE.b64
  cat $OUTPUT_FILE.b64
  echo "===== END $OUTPUT_FILE (Base64) ====="
  '
```
