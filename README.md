# PorkLauncher

Лаунчер сервера PorkLand, based on ArchDev HyLauncher.

## For demo only. Support Hytale developers!

<p align="center">
  <img src="build/appicon.png" alt="PorkLauncher" width="128"/>
</p>

<p align="center">
  <b>Unofficial Hytale Launcher</b><br>
  <i>Неофициальный Hytale лаунчер</i>
</p>
<p align="center">
  <a href="https://github.com/Serezjjja/PorkLauncher/releases"><img alt="GitHub Downloads (all assets, all releases)" src="https://img.shields.io/github/downloads/Serezjjja/PorkLauncher/total?style=flat-square"></a>
  <img src="https://img.shields.io/badge/License-GPL_3.0-yellow?style=flat-square"/>
</p>

---

## Фичи

- Онлайн режим
- Скачивание игры
- Скачивание всех зависимостей
- Уникальные идентификаторы ников (каждый ник уникальный)
- Поддержка Windows и macOS

---

## Установка

Переходим в раздел [releases](https://github.com/Serezjjja/PorkLauncher/releases). <br>
Скачиваем самую [последнюю версию](https://github.com/Serezjjja/PorkLauncher/releases/latest) лаунчера. <br>
Не нужно скачивать `update-helper(.exe)`

---

## Билд

Зависимости
- Golang 1.24+
- NodeJS 24+

### Windows
```
git clone https://github.com/Serezjjja/PorkLauncher.git
cd PorkLauncher
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails build -platform windows/amd64
```
Билд появится в папке `build/bin`

### macOS

macOS builds work on both Intel and Apple Silicon (M1/M2/M3) Macs.

**Build for macOS:**
```bash
wails build -platform darwin/universal
```

**Running the app:**
Since the app is not code-signed, macOS will show security warnings on first launch. To run:
- Right-click the app → "Open" → "Open anyway"
- Or run: `xattr -cr build/bin/HyLauncher.app` (removes quarantine attributes for testing)

**Game launch issues:**
If you get "Hytale.app is corrupted" error when launching the game:
1. The launcher automatically tries to fix permissions - just try launching again
2. If it still fails, manually run:
   ```bash
   xattr -cr ~/Library/Application\ Support/PorkLauncher/shared/games/release/*/Client/Hytale.app
   codesign --force --deep --sign - ~/Library/Application\ Support/PorkLauncher/shared/games/release/*/Client/Hytale.app
   ```
3. Or disable Gatekeeper temporarily: `sudo spctl --master-disable` (re-enable after with `sudo spctl --master-enable`)

**GitHub Actions:**
macOS builds are handled automatically by GitHub Actions. The workflow builds universal binaries that work on both Intel and Apple Silicon Macs.

---

## License

У нас используется лицензия [GPL 3.0](https://choosealicense.com/licenses/gpl-3.0/).<br>
`Permissions of this strong copyleft license are conditioned on making available complete source code of licensed works and modifications, which include larger works using a licensed work, under the same license. Copyright and license notices must be preserved. Contributors provide an express grant of patent rights.` via [choosealicense.com](https://choosealicense.com/licenses)

---

## Credits
- [Hytale-F2P](https://github.com/amiayweb/Hytale-F2P) Online fix method by game patching

---

## Authors

- [@Serezjjja](https://www.github.com/Serezjjja) (Developer)

## Star History

<a href="https://www.star-history.com/#Serezjjja/PorkLauncher&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Serezjjja/PorkLauncher&type=date&theme=dark&legend=top-left" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Serezjjja/PorkLauncher&type=date&legend=top-left" />
   <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=Serezjjja/PorkLauncher&type=date&legend=top-left" />
 </picture>
</a>
