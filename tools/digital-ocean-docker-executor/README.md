## Usage

```bash
zmap -p 80 -o zmap.txt
```

Run executor with 4 droplets, then download the results to local file.

```bash
go run main.go \
    --name http-grab \
    --droplet-public-key-path ./.ssh/executor.pub \
    --droplet-private-key-path ./.ssh/executor \
    --do-token dop_v1_0000000000000000000000000000000000000000000000000000000000000000 \
    --num-droplets 4
```

Run executor with 4 droplets, then upload the results to amazon s3.

```bash
go run main.go \
    --name http-grab \
    --droplet-public-key-path ./.ssh/executor.pub \
    --droplet-private-key-path ./.ssh/executor \
    --do-token dop_v1_0000000000000000000000000000000000000000000000000000000000000000 \
    --s3-access-key=00000000000000000000 \
    --s3-secret-key=0000000000000000000000000000000000000000 \
    --s3-bucket=dode \
    --num-droplets 4

```