#!/bin/bash

# --- 0. 定数・設定 ---
LOCAL_SOCKET="./incus.socket"
INFRA_DIR="./infra"
BACKEND_DIR="./backend"
FRONTEND_DIR="./frontend"

# 色付け用
GREEN='\033[0;32m'
NC='\033[0m' # No Color

echo -e "${GREEN}>>> Initializing Dev Environment...${NC}"

# --- 1. インフラの起動 ---
echo -e "${GREEN}>>> Starting Vagrant VM...${NC}"
cd "$INFRA_DIR" || exit
vagrant up
cd ..

# --- 2. 既存ソケットの掃除 ---
rm -f "$LOCAL_SOCKET"

# --- 3. クリーンアップ関数の定義 ---
cleanup() {
    echo -e "\n${GREEN}>>> Cleaning up and shutting down...${NC}"
    # 各プロセスの停止
    kill $TUNNEL_PID $BACKEND_PID $FRONTEND_PID 2>/dev/null
    
    # VMの停止
    echo "Halt Vagrant VM..."
    cd "$INFRA_DIR" && vagrant halt
    
    # ソケットの削除
    rm -f "$LOCAL_SOCKET"
    echo -e "${GREEN}>>> All processes stopped. See you next time!${NC}"
    exit
}

# Ctrl+C 等を検知した時に cleanup を実行
trap cleanup INT TERM EXIT

# --- 4. SSHトンネル起動 (Background) ---
# Vagrantの秘密鍵パスを取得してトンネルを貼る
echo -e "${GREEN}>>> Establishing SSH Tunnel for Incus Socket...${NC}"
# .vagrant 以下のディレクトリ構成はプロバイダによって変わるため、必要に応じて修正してください
VAGRANT_KEY=".vagrant/machines/default/virtualbox/private_key"

ssh -i "$VAGRANT_KEY" \
    -o StrictHostKeyChecking=no \
    -nNT -L "$LOCAL_SOCKET":/var/lib/incus/unix.socket \
    vagrant@192.168.56.10 &
TUNNEL_PID=$!

# --- 5. バックエンド (Go) 起動 (Background) ---
echo -e "${GREEN}>>> Starting Go Backend (API & OIDC Gate)...${NC}"
export APP_ENV=development
export INCUS_SOCKET="$LOCAL_SOCKET"
export OIDC_CLIENT_ID="my-vm-console-client"
export OIDC_ISSUER="https://idp.ituki.dev" 

cd "$BACKEND_DIR" || exit
go run main.go &
BACKEND_PID=$!
cd ..

# --- 6. フロントエンド (React) 起動 (Foreground) ---
echo -e "${GREEN}>>> Starting React Frontend (Vite)...${NC}"
cd "$FRONTEND_DIR" || exit
# Viteの出力をこのスクリプトの標準出力に流す
npm run dev &
FRONTEND_PID=$!
cd ..

# フロントエンドが動いている間、スクリプトを維持
wait