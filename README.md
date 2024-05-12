# esbuild-plugin-tailwind

An esbuild plugin to use tailwind under a go configuration.

## Install
```sh
go get github.com/iamajoe/esbuild-plugin-tailwind
```

### Dependencies
[tailwindcss cli](https://tailwindcss.com/docs/installation) is required in order to use the plugin.

The plugin will look for `tailwindcss` in 2 places:
- `node_modules` project folder (recursively going through parents until is found)
- on `$PATH`

## Usage
```css
/* input.css */
@tailwind base;
@tailwind components;
@tailwind utilities;
```

```js
// input.js
import "./input.css";
```

```go
import (
    "github.com/evanw/esbuild/pkg/api"
    estailwind "github.com/iamajoe/esbuild-plugin-tailwind"
)

func main() {
	_ = api.Build(api.BuildOptions{
		EntryPoints: []string{"input.js"},
		Outfile:     "output.js",
		Loader: map[string]api.Loader{
			".js":    api.LoaderJS,
			".css":   api.LoaderCSS,
		},
		Plugins:           []api.Plugin{
            // NewTailwindPlugin(minify bool) api.Plugin
            estailwind.NewTailwindPlugin(false),
        },
	})
}
```
