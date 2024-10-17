# ZAPP
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fironpark%2Fzapp.svg?type=shield&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fironpark%2Fzapp?ref=badge_shield&issueType=license)
[![Go Report Card](https://goreportcard.com/badge/github.com/ironpark/zapp)](https://goreportcard.com/report/github.com/ironpark/zapp)
[![codebeat badge](https://codebeat.co/badges/6b004587-036c-4324-bc97-c2e76d58b474)](https://codebeat.co/projects/github-com-ironpark-zapp-main)
![GitHub Repo stars](https://img.shields.io/github/stars/ironpark/zapp)


🌐 [English](README.md) | [한국어](README.ko.md) | [**日本語**](README.ja.md)

![asd](/docs/demo.gif)

**macOSアプリのデプロイメントを簡素化**

`zapp`は、macOSアプリケーションのデプロイメントプロセスを合理化し、自動化するために設計された強力なCLIツールです。依存関係のバンドルからDMG/PKGの作成、コード署名、公証に至るまで、デプロイメントのすべての段階を1つのツールで処理します。

## ✨ 機能

- [x] DMGファイルの作成
- [x] PKGファイルの作成
- [x] コード署名
- [x] 公証 / ステープリング
- [ ] plistの修正（バージョン）
- [x] バイナリ依存関係の自動バンドル
- [ ] GitHub Actionsのサポート

## ⚡️ クイックスタート
#### 🍺 Homebrewを使用
```bash
brew tap ironpark/zapp
brew install zapp
```

#### 🛠️ ソースコードからビルド

```bash
go install github.com/ironpark/zapp@latest
```

## 📖 使用方法
### 🔏 コード署名

> [!TIP]
>
> `--identity`フラグを使用して証明書を選択しない場合、Zappは現在のキーチェーンから利用可能な証明書を自動的に選択します。

```bash
zapp sign --target="path/to/target.(app,dmg,pkg)"
```
```bash
zapp sign --identity="Developer ID Application" --target="path/to/target.(app,dmg,pkg)"
```

### 🏷️ 公証 & ステープリング
> [!NOTE]
>
> 公証コマンドを実行する際、Zappがアプリバンドルパスを受け取ると、自動的にアプリバンドルを圧縮し、公証を試みます。

```bash
zapp notarize --profile="key-chain-profile" --target="path/to/target.(app,dmg,pkg)" --staple
```

```bash
zapp notarize --apple-id="your@email.com" --password="pswd" --team-id="XXXXX" --target="path/to/target.(app,dmg,pkg)" --staple
```

### 🔗 依存関係のバンドル
> [!NOTE]
>
> このプロセスは、アプリケーション実行ファイルの依存関係を検査し、必要なライブラリを `/Contents/Frameworks` 内に含め、スタンドアロン実行を可能にするためにリンクパスを修正します。

```bash
zapp dep --app="path/to/target.app"
```
#### ライブラリを検索する追加パス
```bash
zapp dep --app="path/to/target.app" --libs="/usr/local/lib" --libs="/opt/homebrew/Cellar/ffmpeg/7.0.2/lib"
```
#### 署名 & 公証 & ステープリングを含む
> [!TIP]
>
> `dep`、`dmg`、`pkg`コマンドは、`--sign`、`--notarize`、`--staple`フラグと共に使用できます。
> - `--sign`フラグは、依存関係のバンドル後にアプリバンドルを自動的に署名します。
> - `--notarize`フラグは、署名後にアプリバンドルを自動的に公証します。

```bash
zapp dep --app="path/to/target.app" --sign --notarize --profile "profile" --staple
```

### 💽 DMGファイルの作成

> Zappを使用して、macOSアプリの配布によく使用されるDMGファイルを作成できます。
> アプリバンドルからアイコンを自動的に抽出し、ディスクアイコンを合成し、アプリのドラッグ＆ドロップインストール用のインターフェースを提供することで、DMG作成プロセスを大幅に簡素化します。

```bash
zapp dmg --app="path/to/target.app"
```

```bash
zapp dmg --title="My App" \ 
  --app="path/to/target.app" \
  --icon="path/to/icon.icns" \
  --bg="path/to/background.png" \ 
  --out="MyApp.dmg"
```
#### 署名 & 公証 & ステープリングを含む
> [!TIP]
>
> `dep`、`dmg`、`pkg`コマンドは、`--sign`、`--notarize`、`--staple`フラグと共に使用できます。
> - `--sign`フラグは、依存関係のバンドル後にアプリバンドルを自動的に署名します。
> - `--notarize`フラグは、署名後にアプリバンドルを自動的に公証します。

```bash
zapp dmg --app="path/to/target.app" --sign --notarize --profile "profile" --staple
```

### 📦 PKGファイルの作成

> [!TIP]
>
> `--version`と`--identifier`フラグが設定されていない場合、これらの値は提供されたアプリバンドルのInfo.plistファイルから自動的に取得されます。

#### アプリバンドルからPKGファイルを作成
```bash
zapp pkg --app="path/to/target.app"
```

```bash
zapp pkg --out="MyApp.pkg" --version="1.2.3" --identifier="com.example.myapp" --app="path/to/target.app"
```

#### EULAファイルを含む

複数の言語でのエンドユーザーライセンス契約（EULA）ファイルを含める：

```bash
zapp pkg --eula=en:eula_en.txt,es:eula_es.txt,fr:eula_fr.txt --app="path/to/target.app" 
```
#### 署名 & 公証 & ステープリングを含む
> [!TIP]
>
> `dep`、`dmg`、`pkg`コマンドは、`--sign`、`--notarize`、`--staple`フラグと共に使用できます。
> - `--sign`フラグは、依存関係のバンドル後にアプリバンドルを自動的に署名します。
> - `--notarize`フラグは、署名後にアプリバンドルを自動的に公証します。

```bash
zapp pkg --app="path/to/target.app" --sign --notarize --profile "profile" --staple
```

### 完全な例
以下は、`zapp`を使用して`MyApp.app`の依存関係のバンドル、コード署名、パッケージング、公証、ステープリングを行う方法を示す完全な例です：

```bash
# 依存関係のバンドル
zapp dep --app="MyApp.app"

# コード署名 / 公証 / ステープリング
zapp sign --target="MyApp.app"
zapp notarize --profile="key-chain-profile" --target="MyApp.app" --staple

# pkg/dmgファイルの作成
zapp pkg --app="MyApp.app" --out="MyApp.pkg"
zapp dmg --app="MyApp.app" --out="MyApp.dmg"

# pkg/dmgのコード署名 / 公証 / ステープリング
zapp sign --target="MyApp.app"
zapp sign --target="MyApp.pkg"

zapp notarize --profile="key-chain-profile" --target="MyApp.pkg" --staple
zapp notarize --profile="key-chain-profile" --target="MyApp.dmg" --staple
```
または、簡単な省略コマンドを使用してください
```bash
zapp dep --app="MyApp.app" --sign --notarize --staple

zapp pkg --out="MyApp.pkg" --app="MyApp.app" \ 
  --sign --notarize --profile="key-chain-profile" --staple

zapp dmg --out="MyApp.dmg" --app="MyApp.app" \
  --sign --notarize --profile="key-chain-profile" --staple
```

## ライセンス
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fironpark%2Fzapp.svg?type=large&issueType=license)](https://app.fossa.com/projects/git%2Bgithub.com%2Fironpark%2Fzapp?ref=badge_large&issueType=license)

Zappは[MITライセンス](LICENSE)の下で公開されています。

## サポート

問題が発生したり質問がある場合は、[GitHubイシュートラッカー](https://github.com/ironpark/zapp/issues)にイシューを作成してください。