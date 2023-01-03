# redis clear keys

use redis pipeline to clear keys that never expire

## How to install?

1. Install Go version at least **1.17**
    * See: [Go installation instructions](https://go.dev/doc/install)
2. Run command:
   ```bash
   go install github.com/xyctruth/redis-clear-keys@latest
   ```
3. Add following line in your `.bashrc`/`.zshrc` file:
   ```bash
   export PATH="$PATH:$HOME/go/bin"
   ```

## How to use?

```bash
redis-clear-keys --redis-host={redis-host} --redis-password={redis-password}
```
