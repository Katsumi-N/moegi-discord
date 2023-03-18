# moegi-discord
## About
Goで書かれたdiscord botです．
ConohaVPSでマインクラフト用のサーバーを立てており、そのサーバーの起動、終了等をdiscordのbotを通じて行います

## command
- !conoha start -> サーバーを起動します(メモリタイプを1gb->4gbに変更した上でサーバーを起動します)
- !conoha stop -> サーバーを終了します(シャットダウンした上でメモリタイプを4gb->1gbに変更します)
- !conoha reboot -> サーバーを再起動します
- !vote　(タイトル)　(選択肢1) (選択肢2) ... -> 投票を作成します(--crirona を末尾につけると3分後に投票していない人にリマインドします)
