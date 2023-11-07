# Jitterbit Extractor

 A desktop app for Jitterbit script import from local project files made with [Wails](https://wails.io/) and [Svelte](https://github.com/BillBuilt/wails-vite-svelte-tailwind-template). Imports script code from selected environment into `.jb` files preserving the original directory structure.

## Building
[Install Wails](https://wails.io/docs/gettingstarted/installation) (requires Go and Node), then check it with:
```
wails doctor
```

After compilation, the build can be found in `build/bin`:
 ```
 wails build
 ```

## Debugging

Create an application log with:
```
JitterbitExtractor.exe > log.txt 2>&1
```
