#!/bin/bash

echo "Please enter or create a working directory:"
read -r WORK_DIR

mkdir -p "$WORK_DIR" # Creates the working directory if non existing

cd "$WORK_DIR"

curl -LO https://github.com/89luca89/distrobox/archive/refs/tags/1.5.0.2.tar.gz
tar -xzf 1.5.0.2.tar.gz

curl -LO https://github.com/Vanilla-OS/apx/releases/download/continuous/apx.tar.gz
tar -xzf apx.tar.gz
mv apx "$HOME/.local/bin/apx2"
chmod +x "$HOME/.local/bin/apx2"

mkdir -p "$HOME/.config/apx"
echo '{
  "apxPath": "'"$HOME/.local/share/apx/"'",
  "distroboxpath": "'"$WORK_DIR/distrobox-1.5.0.2/distrobox"'",
  "storageDriver": "btrfs"
}' > "$HOME/.config/apx/apx.json"

git clone https://github.com/Vanilla-OS/vanilla-apx-configs.git "$WORK_DIR/vanilla-apx-configs"
mv "$WORK_DIR/vanilla-apx-configs/stacks" "$HOME/.local/share/apx/"
mv "$WORK_DIR/vanilla-apx-configs/package-managers" "$HOME/.local/share/apx/"

echo "Installation completed. You can now use Apx v2 by running 'apx2'."
apx2 --version
