## Usage

```bash
zmap -p 80 -o zmap.txt
```

```
go run main.go \
    --name http-grab \
    --droplet-public-key-path ./.ssh/executor.pub \
    --droplet-private-key-path ./.ssh/executor \
    --do-token dop_v1_0000000000000000000000000000000000000000000000000000000000000000 \
    --num-droplets 4
```