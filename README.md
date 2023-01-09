# redis-clear-keys

Use the redis pipeline parallel batch to remove the never expired key

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

### Argument

```bash
redis-clear-keys --redis-host={redis-host} --redis-password={redis-password}

redis-clear-keys -h 
Usage of redis-clear-keys:
  -batch-number int
        number of batches (default 1000)
  -parallel-number int
        number of parallel processing (default 100)
  -redis-host string
        redis host
  -redis-password string
        redis password
```
